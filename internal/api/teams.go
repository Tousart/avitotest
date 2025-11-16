package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tousart/avitotest/internal/api/helpers"
	"github.com/tousart/avitotest/internal/api/types"
	"github.com/tousart/avitotest/internal/codes"
	"github.com/tousart/avitotest/internal/usecase"
)

type Teams struct {
	teamsService usecase.TeamsService
}

func CreateTeamsAPI(teamsService usecase.TeamsService) *Teams {
	return &Teams{
		teamsService: teamsService,
	}
}

func (t *Teams) teamAddHandler(w http.ResponseWriter, r *http.Request) {
	team, err := types.CreateTeamAddRequest(r)
	if err != nil {
		helpers.WriteAPIError(w, http.StatusBadRequest, codes.ErrBadRequet, err.Error())
		return
	}

	errResp := t.teamsService.TeamAdd(r.Context(), team)
	if errResp != nil {
		status := helpers.GetStatusError(errResp.Code)
		helpers.WriteAPIError(w, status, errResp.Code, errResp.Message)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(team)
}

func (t *Teams) teamGetHandler(w http.ResponseWriter, r *http.Request) {
	team, err := types.CreateTeamGetRequest(r)
	if err != nil {
		helpers.WriteAPIError(w, http.StatusBadRequest, codes.ErrBadRequet, err.Error())
		return
	}

	errResp := t.teamsService.TeamGet(r.Context(), team)
	if errResp != nil {
		status := helpers.GetStatusError(errResp.Code)
		helpers.WriteAPIError(w, status, errResp.Code, errResp.Message)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(team)
}

func (t *Teams) WithTeamsHandlers(r chi.Router) {
	r.Route("/team", func(r chi.Router) {
		r.Post("/add", t.teamAddHandler)
		r.Get("/get", t.teamGetHandler)
	})
}
