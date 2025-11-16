package service

import (
	"context"

	"github.com/tousart/avitotest/internal/models"
	"github.com/tousart/avitotest/internal/repository"
)

const (
	TimeFormat      = "2006-01-02T15:04:05Z"
	DefaultPRStatus = "OPEN"
)

type PullRequestsService struct {
	repo repository.PullRequestsRepository
}

func NewPullRequestsService(repo repository.PullRequestsRepository) *PullRequestsService {
	return &PullRequestsService{
		repo: repo,
	}
}

func (ps *PullRequestsService) PullRequestCreate(ctx context.Context, pullRequest *models.PullRequest) *models.ErrorResponse {
	pullRequest.Status = DefaultPRStatus

	err, createdAt, reviewers := ps.repo.PullRequestCreate(ctx, pullRequest)
	if err != nil {
		return err
	}

	pullRequest.AssignedReviewers = reviewers
	pullRequest.CreatedAt = (*createdAt).Format(TimeFormat)

	return nil
}

func (ps *PullRequestsService) PullRequestMerge(ctx context.Context, pullRequest *models.PullRequest) *models.ErrorResponse {
	err, createdAt, mergedAt, pullRequestName, authorID, status := ps.repo.PullRequestMerge(ctx, pullRequest)
	if err != nil {
		return err
	}

	pullRequest.PullRequestName = pullRequestName
	pullRequest.AuthorID = authorID
	pullRequest.Status = status
	pullRequest.CreatedAt = (*createdAt).Format(TimeFormat)
	pullRequest.MergedAt = (*mergedAt).Format(TimeFormat)

	return nil
}

func (ps *PullRequestsService) PullRequestReassign(ctx context.Context, pullRequest *models.PullRequest, oldUserID string) *models.ErrorResponse {
	err, pullRequestName, authorID, status, reviewers := ps.repo.PullRequestReassign(ctx, pullRequest, oldUserID)
	if err != nil {
		return err
	}

	pullRequest.PullRequestName = pullRequestName
	pullRequest.AuthorID = authorID
	pullRequest.Status = status
	pullRequest.AssignedReviewers = reviewers

	return nil
}
