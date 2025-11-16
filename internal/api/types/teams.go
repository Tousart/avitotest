package types

import (
	"encoding/json"
	"net/http"

	"errors"

	"github.com/tousart/avitotest/internal/models"
)

func CreateTeamAddRequest(r *http.Request) (*models.Team, error) {
	var request models.Team

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	if request.TeamName == "" {
		return nil, errors.New("team name is required")
	}

	if len(request.Members) == 0 {
		return nil, errors.New("members are required")
	}

	return &request, nil
}

func CreateTeamGetRequest(r *http.Request) (*models.Team, error) {
	var request models.Team
	request.TeamName = r.URL.Query().Get("team_name")

	if request.TeamName == "" {
		return nil, errors.New("team name is required")
	}

	return &request, nil
}
