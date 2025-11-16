package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
	"github.com/tousart/avitotest/internal/codes"
	"github.com/tousart/avitotest/internal/models"
	"github.com/tousart/avitotest/pkg"
)

type TeamsRepository struct {
	db *sql.DB
}

func NewTeamsRepository(addressToConnectToPSQL string) (*TeamsRepository, error) {
	db, err := pkg.ConnectToPSQL(addressToConnectToPSQL)
	if err != nil {
		log.Printf("failed to connect to db: %v\n", err)
		return nil, fmt.Errorf("repository: postgres: NewTeamsRepository: %v", err)
	}

	return &TeamsRepository{db: db}, nil
}

func (tr *TeamsRepository) TeamAdd(ctx context.Context, team *models.Team) *models.ErrorResponse {
	// Начинаем транзакцию

	tx, err := tr.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("repository: postgres: TeamAdd: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}
	}

	// Добавление команды

	queryCreateTeam := "INSERT INTO teams (team_name) values ($1);"
	_, err = tx.ExecContext(ctx, queryCreateTeam, team.TeamName)
	if err != nil {
		tx.Rollback()
		log.Printf("repository: postgres: TeamAdd: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrTeamExists,
			Message: "team_name already exists",
		}
	}

	usersID := make([]string, len(team.Members))
	for i, teamMember := range team.Members {
		usersID[i] = teamMember.UserID
	}

	// Поиск существующих пользователей

	querySelectExistsUsers := "SELECT user_id FROM users WHERE user_id IN (SELECT * FROM unnest($1::varchar[]));"
	rowsExists, err := tx.QueryContext(ctx, querySelectExistsUsers, pq.Array(usersID))
	if err != nil && err != sql.ErrNoRows {
		log.Printf("repository: postgres: TeamAdd: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}
	}

	existsUsers := make([]string, 0)
	existsUsersMap := make(map[string]struct{})

	for rowsExists.Next() {
		var userID string

		if err := rowsExists.Scan(&userID); err != nil {
			return &models.ErrorResponse{
				Code:    codes.ErrInternal,
				Message: "internal error",
			}
		}

		existsUsers = append(existsUsers, userID)
		existsUsersMap[userID] = struct{}{}
	}
	rowsExists.Close()

	// Обновление команды у пользователей, которые существуют

	queryUpdateUsersTeam := "UPDATE users SET team_name = $1 WHERE user_id IN (SELECT * FROM unnest($2::varchar[]));"
	_, err = tx.ExecContext(ctx, queryUpdateUsersTeam, team.TeamName, pq.Array(existsUsers))
	if err != nil {
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}
	}

	// Добавление пользователей, которых не существовало

	newUsersID := make([]string, 0)
	newUsersUsername := make([]string, 0)
	newUsersIsActive := make([]bool, 0)

	for _, teamMember := range team.Members {
		if _, ok := existsUsersMap[teamMember.UserID]; !ok {
			newUsersID = append(newUsersID, teamMember.UserID)
			newUsersUsername = append(newUsersUsername, teamMember.Username)
			newUsersIsActive = append(newUsersIsActive, teamMember.IsActive)
		}
	}

	if len(newUsersID) > 0 {
		queryInsertNewUsers := `
		INSERT INTO users (user_id, username, team_name, is_active) 
		SELECT u_id, u_name, $3, u_active
		FROM unnest($1::varchar[], $2::varchar[], $4::boolean[]) AS t(u_id, u_name, u_active);
		`
		_, err = tx.ExecContext(ctx, queryInsertNewUsers, pq.Array(newUsersID), pq.Array(newUsersUsername), team.TeamName, pq.Array(newUsersIsActive))
		if err != nil {
			return &models.ErrorResponse{
				Code:    codes.ErrInternal,
				Message: "internal error",
			}
		}
	}

	// Коммит

	if err := tx.Commit(); err != nil {
		log.Printf("repository: postgres: TeamAdd: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}
	}

	return nil
}

func (tr *TeamsRepository) TeamGet(ctx context.Context, team *models.Team) (*models.ErrorResponse, []models.TeamMember) {
	// Проверка: существует ли команда

	var exists bool
	queryTeamExists := "SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1);"
	err := tr.db.QueryRowContext(ctx, queryTeamExists, team.TeamName).Scan(&exists)
	if err != nil {
		log.Printf("repository: postgres: TeamGet: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}, nil
	} else if !exists {
		return &models.ErrorResponse{
			Code:    codes.ErrNotFound,
			Message: "team not found",
		}, nil
	}

	// Получение пользователей по названию команды

	queryGetMembers := "SELECT user_id, username, is_active FROM users WHERE team_name = $1;"
	rows, err := tr.db.QueryContext(ctx, queryGetMembers, team.TeamName)
	if err != nil {
		log.Printf("repository: postgres: TeamGet: %v\n", err)
		return &models.ErrorResponse{
			Code:    codes.ErrInternal,
			Message: "internal error",
		}, nil
	}
	defer rows.Close()

	members := make([]models.TeamMember, 0)
	for rows.Next() {
		var member models.TeamMember

		if err := rows.Scan(&member.UserID, &member.Username, &member.IsActive); err != nil {
			return &models.ErrorResponse{
				Code:    codes.ErrInternal,
				Message: "internal error",
			}, nil
		}

		members = append(members, member)
	}

	return nil, members
}
