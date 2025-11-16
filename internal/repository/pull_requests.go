package repository

import (
	"context"
	"time"

	"github.com/tousart/avitotest/internal/models"
)

type PullRequestsRepository interface {
	PullRequestCreate(ctx context.Context, pullRequest *models.PullRequest) (*models.ErrorResponse, *time.Time, []string)
	PullRequestMerge(ctx context.Context, pullRequest *models.PullRequest) (*models.ErrorResponse, *time.Time, *time.Time, string, string, string)
	PullRequestReassign(ctx context.Context, pullRequest *models.PullRequest, oldUserID string) (*models.ErrorResponse, string, string, string, []string)
}
