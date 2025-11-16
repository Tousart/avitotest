package usecase

import (
	"context"

	"github.com/tousart/avitotest/internal/models"
)

type UsersService interface {
	SetIsActive(ctx context.Context, user *models.User) *models.ErrorResponse
	GetReview(ctx context.Context, pullRequests *[]models.PullRequestShort, userID string) *models.ErrorResponse
	GetActivity(ctx context.Context, usersActivity *[]models.UserActivity) *models.ErrorResponse
}
