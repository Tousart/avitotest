package repository

import (
	"context"

	"github.com/tousart/avitotest/internal/models"
)

type TeamsRepository interface {
	TeamAdd(ctx context.Context, team *models.Team) *models.ErrorResponse
	TeamGet(ctx context.Context, team *models.Team) (*models.ErrorResponse, []models.TeamMember)
}
