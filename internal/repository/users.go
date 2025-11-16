package repository

import (
	"context"

	"github.com/tousart/avitotest/internal/models"
)

type UsersRepository interface {
	SetIsActive(ctx context.Context, user *models.User) (*models.ErrorResponse, string, string)
	GetReview(ctx context.Context, userID string) (*models.ErrorResponse, []models.PullRequestShort)
	GetActivity(ctx context.Context) (*models.ErrorResponse, []models.UserActivity)
}
