package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/tousart/avitotest/internal/codes"
	"github.com/tousart/avitotest/internal/models"
	"github.com/tousart/avitotest/pkg"
)

type UsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(addressToConnectToPSQL string) (*UsersRepository, error) {
	db, err := pkg.ConnectToPSQL(addressToConnectToPSQL)
	if err != nil {
		log.Printf("failed to connect to db: %v\n", err)
		return nil, fmt.Errorf("repository: postgres: NewUsersRepository: %v", err)
	}

	return &UsersRepository{db: db}, nil
}

func (ur *UsersRepository) SetIsActive(ctx context.Context, user *models.User) (*models.ErrorResponse, string, string) {
	var (
		username string
		teamName string
	)

	querySetIsActive := "UPDATE users SET is_active = $1 WHERE user_id = $2 RETURNING username, team_name;"
	err := ur.db.QueryRowContext(ctx, querySetIsActive, user.IsActive, user.UserID).Scan(&username, &teamName)
	if err == sql.ErrNoRows {
		return &models.ErrorResponse{
			Code:    codes.ErrNotFound,
			Message: "user not found",
		}, "", ""
	} else if err != nil {
		log.Printf("repository: postgres: SetIsActive: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}, "", ""
	}
	return nil, username, teamName
}

func (ur *UsersRepository) GetReview(ctx context.Context, userID string) (*models.ErrorResponse, []models.PullRequestShort) {
	var exists bool
	queryExists := "SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1);"
	err := ur.db.QueryRowContext(ctx, queryExists, userID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("repository: postgres: GetReview: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}, nil
	} else if !exists {
		return &models.ErrorResponse{
			Code:    codes.ErrNotFound,
			Message: "user not found",
		}, nil
	}

	queryGetPullRequests := `
	SELECT
		p.pull_request_id,
		pr.pull_request_name,
		pr.author_id,
		pr.status
	FROM pr_reviewers p
	JOIN pull_requests pr ON pr.pull_request_id = p.pull_request_id
	WHERE p.user_id = $1;
	`
	rows, err := ur.db.QueryContext(ctx, queryGetPullRequests, userID)
	if err != nil {
		log.Printf("repository: postgres: GetReview: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}, nil
	}
	defer rows.Close()

	pullRequests := make([]models.PullRequestShort, 0)

	for rows.Next() {
		var pullRequestShort models.PullRequestShort

		if err := rows.Scan(
			&pullRequestShort.PullRequestID, &pullRequestShort.PullRequestName, &pullRequestShort.AuthorID, &pullRequestShort.Status); err != nil {
			log.Printf("repository: postgres: GetReview: %v\n", err)
			return &models.ErrorResponse{
				Code:    codes.ErrInternal,
				Message: "internal error",
			}, nil
		}

		pullRequests = append(pullRequests, pullRequestShort)
	}

	return nil, pullRequests
}

func (ur *UsersRepository) GetActivity(ctx context.Context) (*models.ErrorResponse, []models.UserActivity) {
	queryGetActivity := `
	SELECT 
		u.user_id, 
		u.username, 
		COUNT(p.pull_request_id) as pull_requests, 
		COUNT(CASE WHEN pr.status = 'MERGED' THEN 1 END) AS merged_pr, 
		COUNT(CASE WHEN pr.status = 'OPEN' THEN 1 END) as open_pr 
	FROM users u 
	JOIN pr_reviewers p USING(user_id) 
	JOIN pull_requests pr USING(pull_request_id)
	GROUP BY u.username, u.user_id
	ORDER BY pull_requests DESC;
	`
	rows, err := ur.db.QueryContext(ctx, queryGetActivity)
	if err != nil {
		log.Printf("repository: postgres: GetActivity: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}, nil
	}
	defer rows.Close()

	activity := make([]models.UserActivity, 0)

	for rows.Next() {
		var userActivity models.UserActivity

		if err := rows.Scan(&userActivity.UserID, &userActivity.Username, &userActivity.PullRequests, &userActivity.MergedPR, &userActivity.OpenPR); err != nil {
			log.Printf("repository: postgres: GetActivity: %v\n", err)
			return &models.ErrorResponse{
				Code:    codes.ErrInternal,
				Message: "internal error",
			}, nil
		}

		activity = append(activity, userActivity)
	}

	return nil, activity
}
