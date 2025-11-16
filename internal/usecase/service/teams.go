package service

import (
	"context"

	"github.com/tousart/avitotest/internal/models"
	"github.com/tousart/avitotest/internal/repository"
)

type TeamsService struct {
	repo repository.TeamsRepository
}

func NewTeamsService(repo repository.TeamsRepository) *TeamsService {
	return &TeamsService{
		repo: repo,
	}
}

func (ts *TeamsService) TeamAdd(ctx context.Context, team *models.Team) *models.ErrorResponse {
	err := ts.repo.TeamAdd(ctx, team)
	if err != nil {
		return err
	}
	return nil
}

func (ts *TeamsService) TeamGet(ctx context.Context, team *models.Team) *models.ErrorResponse {
	err, members := ts.repo.TeamGet(ctx, team)
	if err != nil {
		return err
	}

	team.Members = members

	return nil
}
