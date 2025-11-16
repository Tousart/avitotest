package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	"github.com/tousart/avitotest/internal/codes"
	"github.com/tousart/avitotest/internal/models"
	"github.com/tousart/avitotest/pkg"
)

const StatusMerged = "MERGED"

type PullRequestsRepository struct {
	db *sql.DB
}

func NewPullRequestsRepository(addressToConnectToPSQL string) (*PullRequestsRepository, error) {
	db, err := pkg.ConnectToPSQL(addressToConnectToPSQL)
	if err != nil {
		log.Printf("failed to connect to db: %v\n", err)
		return nil, fmt.Errorf("repository: postgres: NewPullRequestsRepository: %v", err)
	}

	return &PullRequestsRepository{db: db}, nil
}

func (pr *PullRequestsRepository) PullRequestCreate(ctx context.Context, pullRequest *models.PullRequest) (*models.ErrorResponse, *time.Time, []string) {
	// Начинаем транзакцию

	tx, err := pr.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("repository: postgres: PullRequestCreate: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}, nil, nil
	}

	// Проверка на существование пулл реквеста

	var existsPullRequest bool
	queryPullRequestExists := "SELECT EXISTS (SELECT 1 FROM pull_requests WHERE pull_request_id = $1);"
	err = tx.QueryRowContext(ctx, queryPullRequestExists, pullRequest.PullRequestID).Scan(&existsPullRequest)
	if err != nil {
		tx.Rollback()
		log.Printf("repository: postgres: PullRequestCreate: queryPullRequestExists: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}, nil, nil
	} else if existsPullRequest {
		tx.Rollback()
		return &models.ErrorResponse{
			Code:    codes.ErrPRExists,
			Message: "pull request already exists",
		}, nil, nil
	}

	// Проверка на существование автора

	var authorsTeam string
	queryAuthorExists := "SELECT team_name FROM users WHERE user_id = $1;"
	err = tx.QueryRowContext(ctx, queryAuthorExists, pullRequest.AuthorID).Scan(&authorsTeam)
	if err == sql.ErrNoRows {
		tx.Rollback()
		return &models.ErrorResponse{
			Code:    codes.ErrNotFound,
			Message: "author not found",
		}, nil, nil
	} else if err != nil {
		tx.Rollback()
		log.Printf("repository: postgres: PullRequestCreate: queryAuthorExists: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}, nil, nil
	}

	// Добавление пулл реквеста в таблицу (и получение created_at)

	var createdAt time.Time
	queryInsertPR := `
	INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status)
	VALUES ($1, $2, $3, $4)
	RETURNING created_at;
	`
	err = tx.QueryRowContext(ctx, queryInsertPR,
		pullRequest.PullRequestID, pullRequest.PullRequestName, pullRequest.AuthorID, pullRequest.Status).Scan(&createdAt)
	if err != nil {
		tx.Rollback()
		log.Printf("repository: postgres: PullRequestCreate: queryInsertPR: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}, nil, nil
	}

	// Поиск свободных ревьюеров

	queryAssignReviewers := `
	SELECT user_id FROM users 
	WHERE team_name = $1 AND is_active = true AND user_id <> $2
	ORDER BY RANDOM()
	LIMIT 2;
	`
	rows, err := tx.QueryContext(ctx, queryAssignReviewers, authorsTeam, pullRequest.AuthorID)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		log.Printf("repository: postgres: PullRequestCreate: queryAssignReviewers: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}, nil, nil
	}

	reviewers := make([]string, 0)
	for rows.Next() {
		var userID string

		if err := rows.Scan(&userID); err != nil {
			tx.Rollback()
			log.Printf("repository: postgres: PullRequestCreate: user_id: %v\n", err)
			return &models.ErrorResponse{
				Code:    codes.ErrInternal,
				Message: "internal error",
			}, nil, nil
		}

		reviewers = append(reviewers, userID)
	}
	rows.Close()

	// Добавление в pr_reviewers значений pull_request_id - user_id

	queryInsertReviewer := "INSERT INTO pr_reviewers (pull_request_id, user_id) SELECT $1, unnest($2::varchar[]);"
	_, err = tx.ExecContext(ctx, queryInsertReviewer, pullRequest.PullRequestID, pq.Array(reviewers))
	if err != nil {
		tx.Rollback()
		log.Printf("repository: postgres: PullRequestCreate: queryInsertReviewer: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}, nil, nil
	}

	// Коммит

	if err := tx.Commit(); err != nil {
		log.Printf("repository: postgres: PullRequestCreate: commit: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}, nil, nil
	}

	return nil, &createdAt, reviewers
}

func (pr *PullRequestsRepository) PullRequestMerge(ctx context.Context, pullRequest *models.PullRequest) (*models.ErrorResponse, *time.Time, *time.Time, string, string, string) {
	// Начинаем транзакцию

	tx, err := pr.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("repository: postgres: PullRequestMerge: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}, nil, nil, "", "", ""
	}

	// Проверяем, существует ли пулл реквест, и, если существует, берем его данные

	var (
		createdAt       time.Time
		mergedAt        sql.NullTime
		pullRequestName string
		authorID        string
		status          string
	)

	queryExistsPR := "SELECT pull_request_name, author_id, status, created_at, merged_at FROM pull_requests WHERE pull_request_id = $1;"
	err = tx.QueryRowContext(ctx, queryExistsPR, pullRequest.PullRequestID).Scan(
		&pullRequestName, &authorID, &status, &createdAt, &mergedAt)
	if err == sql.ErrNoRows {
		tx.Rollback()
		return &models.ErrorResponse{
			Code:    codes.ErrNotFound,
			Message: "pull request not found",
		}, nil, nil, "", "", ""
	} else if err != nil {
		tx.Rollback()
		log.Printf("repository: postgres: PullRequestMerge: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}, nil, nil, "", "", ""
	}

	// Если пулл реквест уже MERGED, то просто возвращаем объект пулл реквеста (как раз это условие реализует идемпотентность: status и merged_at изменяются только раз)

	if pullRequest.Status == StatusMerged {
		return nil, &createdAt, &mergedAt.Time, pullRequestName, authorID, status
	}

	// Если статус не MERGED (а OPEN), то обновляем статус и фиксируем время (merged_at)

	queryUpdateStatus := "UPDATE pull_requests SET status = 'MERGED', merged_at = NOW() WHERE pull_request_id = $1 RETURNING status, merged_at;"
	err = tx.QueryRowContext(ctx, queryUpdateStatus, pullRequest.PullRequestID).Scan(&status, &mergedAt)
	if err != nil {
		tx.Rollback()
		log.Printf("repository: postgres: PullRequestMerge: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}, nil, nil, "", "", ""
	}

	// Коммит

	if err := tx.Commit(); err != nil {
		log.Printf("repository: postgres: PullRequestMerge: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}, nil, nil, "", "", ""
	}

	return nil, &createdAt, &mergedAt.Time, pullRequestName, authorID, status
}

func (pr *PullRequestsRepository) PullRequestReassign(ctx context.Context, pullRequest *models.PullRequest, oldUserID string) (*models.ErrorResponse, string, string, string, []string) {
	// Начинаем транзакцию

	tx, err := pr.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("repository: postgres: PullRequestReassign: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}, "", "", "", nil
	}

	// Проверка на существование пулл реквеста
	// Проверка на существование old_user_id
	// Проверка: является ли old_user_id ревьюером этого пулл реквеста
	// Проверка статуса пулл реквеста

	var (
		teamName        string
		authorID        string
		status          string
		isReviewer      bool
		pullRequestName string
	)

	queryCheck := `
	SELECT
		u.team_name,
		pr.author_id,
		pr.status,
		pr.pull_request_name,
		EXISTS(
			SELECT 1 FROM pr_reviewers r
			WHERE r.pull_request_id = $1 AND r.user_id = $2
		) AS is_reviewer
	FROM pull_requests pr
	JOIN users u ON u.user_id = $2
	WHERE pr.pull_request_id = $1;
	`

	err = tx.QueryRowContext(ctx, queryCheck, pullRequest.PullRequestID, oldUserID).Scan(&teamName, &authorID, &status, &pullRequestName, &isReviewer)
	if err == sql.ErrNoRows {
		return &models.ErrorResponse{
			Code:    codes.ErrNotFound,
			Message: "pull request not found",
		}, "", "", "", nil
	} else if err != nil {
		log.Printf("repository: postgres: PullRequestReassign: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}, "", "", "", nil
	}

	if teamName == "" {
		return &models.ErrorResponse{
			Code:    codes.ErrNotFound,
			Message: "old user not found",
		}, "", "", "", nil
	}

	if !isReviewer {
		return &models.ErrorResponse{
			Code:    codes.ErrNotAssigned,
			Message: "old user is not a reviewer of this pull request",
		}, "", "", "", nil
	}

	if status == StatusMerged {
		return &models.ErrorResponse{
			Code:    codes.ErrPRMerged,
			Message: "pull request is merged",
		}, "", "", "", nil
	}

	// Получение нового кандидата

	var newUserID string
	queryNewCandidate := `
	SELECT user_id FROM users 
	WHERE 
		is_active = true AND 
		user_id <> $1 AND 
		team_name = $2 AND 
		user_id NOT IN (SELECT user_id FROM pr_reviewers WHERE pull_request_id = $3) 
	ORDER BY RANDOM() 
	LIMIT 1;
	`
	err = tx.QueryRowContext(ctx, queryNewCandidate, authorID, teamName, pullRequest.PullRequestID).Scan(&newUserID)
	if err == sql.ErrNoRows {
		return &models.ErrorResponse{
			Code:    codes.ErrNoCandidate,
			Message: "no available candidates",
		}, "", "", "", nil
	} else if err != nil {
		log.Printf("repository: postgres: PullRequestReassign: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}, "", "", "", nil
	}

	// Заменяем старого кандидата на нового в этом пулл реквесте

	queryUpdateCandidate := "UPDATE pr_reviewers SET user_id = $1 WHERE user_id = $2 AND pull_request_id = $3;"
	_, err = tx.ExecContext(ctx, queryUpdateCandidate, newUserID, oldUserID, pullRequest.PullRequestID)
	if err != nil {
		log.Printf("repository: postgres: PullRequestReassign: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}, "", "", "", nil
	}

	// Получение обновленного набора ревьюеров

	queryGetReviewers := "SELECT user_id FROM pr_reviewers WHERE pull_request_id = $1;"
	rows, err := tx.QueryContext(ctx, queryGetReviewers, pullRequest.PullRequestID)
	if err != nil {
		log.Printf("repository: postgres: PullRequestReassign: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}, "", "", "", nil
	}
	defer rows.Close()

	reviewers := make([]string, 0)
	for rows.Next() {
		var userID string

		if err := rows.Scan(&userID); err != nil {
			log.Printf("repository: postgres: PullRequestReassign: %v\n", err)
			return &models.ErrorResponse{
				Code:    codes.ErrInternal,
				Message: "internal error",
			}, "", "", "", nil
		}

		reviewers = append(reviewers, userID)
	}

	// Коммит

	if err := tx.Commit(); err != nil {
		log.Printf("repository: postgres: PullRequestReassign: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}, "", "", "", nil
	}

	return nil, pullRequestName, authorID, status, reviewers
}
