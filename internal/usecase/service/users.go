package service

import (
	"context"

	"github.com/tousart/avitotest/internal/models"
	"github.com/tousart/avitotest/internal/repository"
)

type UsersService struct {
	repo repository.UsersRepository
}

func NewUsersService(repo repository.UsersRepository) *UsersService {
	return &UsersService{
		repo: repo,
	}
}

func (us *UsersService) SetIsActive(ctx context.Context, user *models.User) *models.ErrorResponse {
	err, username, teamName := us.repo.SetIsActive(ctx, user)
	if err != nil {
		return err
	}

	user.Username = username
	user.TeamName = teamName

	return nil
}

func (us *UsersService) GetReview(ctx context.Context, pullRequests *[]models.PullRequestShort, userID string) *models.ErrorResponse {
	err, PRs := us.repo.GetReview(ctx, userID)
	if err != nil {
		return err
	}

	(*pullRequests) = PRs

	return nil
}

func (us *UsersService) GetActivity(ctx context.Context, usersActivity *[]models.UserActivity) *models.ErrorResponse {
	err, activity := us.repo.GetActivity(ctx)
	if err != nil {
		return err
	}

	(*usersActivity) = activity

	return nil
}
