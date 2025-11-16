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

type PullRequests struct {
	pullRequestsService usecase.PullRequestsService
}

func CreatePullRequestsAPI(pullRequestsService usecase.PullRequestsService) *PullRequests {
	return &PullRequests{
		pullRequestsService: pullRequestsService,
	}
}

func (pr *PullRequests) pullRequestCreateHandler(w http.ResponseWriter, r *http.Request) {
	pullRequest, err := types.CreatePullRequestCreateRequest(r)
	if err != nil {
		helpers.WriteAPIError(w, http.StatusBadRequest, codes.ErrBadRequet, "bad request")
		return
	}

	errResp := pr.pullRequestsService.PullRequestCreate(r.Context(), pullRequest)
	if errResp != nil {
		status := helpers.GetStatusError(errResp.Code)
		helpers.WriteAPIError(w, status, errResp.Code, errResp.Message)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pullRequest)
}

func (pr *PullRequests) pullRequestMergeHandler(w http.ResponseWriter, r *http.Request) {
	pullRequest, err := types.CreatePullRequestMergeRequest(r)
	if err != nil {
		helpers.WriteAPIError(w, http.StatusBadRequest, codes.ErrBadRequet, "bad request")
		return
	}

	errResp := pr.pullRequestsService.PullRequestMerge(r.Context(), pullRequest)
	if errResp != nil {
		status := helpers.GetStatusError(errResp.Code)
		helpers.WriteAPIError(w, status, errResp.Code, errResp.Message)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pullRequest)
}

func (pr *PullRequests) pullRequestReassignHandler(w http.ResponseWriter, r *http.Request) {
	pullRequest, oldUserID, err := types.CreatePullRequestReassign(r)
	if err != nil {
		helpers.WriteAPIError(w, http.StatusBadRequest, codes.ErrBadRequet, "bad request")
		return
	}

	errResp := pr.pullRequestsService.PullRequestReassign(r.Context(), pullRequest, oldUserID)
	if errResp != nil {
		status := helpers.GetStatusError(errResp.Code)
		helpers.WriteAPIError(w, status, errResp.Code, errResp.Message)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pullRequest)
}

func (pr *PullRequests) WithPullRequestsHandlers(r chi.Router) {
	r.Route("/pullRequest", func(r chi.Router) {
		r.Post("/create", pr.pullRequestCreateHandler)
		r.Post("/merge", pr.pullRequestMergeHandler)
		r.Post("/reassign", pr.pullRequestReassignHandler)
	})
}
