package usecase

import (
	"context"

	"github.com/tousart/avitotest/internal/models"
)

type PullRequestsService interface {
	PullRequestCreate(ctx context.Context, pullRequest *models.PullRequest) *models.ErrorResponse
	PullRequestMerge(ctx context.Context, pullRequest *models.PullRequest) *models.ErrorResponse
	PullRequestReassign(ctx context.Context, pullRequest *models.PullRequest, oldUserID string) *models.ErrorResponse
}
