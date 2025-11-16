package usecase

import (
	"context"

	"github.com/tousart/avitotest/internal/models"
)

type TeamsService interface {
	TeamAdd(ctx context.Context, team *models.Team) *models.ErrorResponse
	TeamGet(ctx context.Context, team *models.Team) *models.ErrorResponse
}
