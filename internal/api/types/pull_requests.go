package types

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/tousart/avitotest/internal/models"
)

func CreatePullRequestCreateRequest(r *http.Request) (*models.PullRequest, error) {
	var request models.PullRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	if request.PullRequestID == "" {
		return nil, errors.New("pull request id is required")
	}

	if request.PullRequestName == "" {
		return nil, errors.New("pull request name is required")
	}

	if request.AuthorID == "" {
		return nil, errors.New("author id is required")
	}

	return &request, nil
}

func CreatePullRequestMergeRequest(r *http.Request) (*models.PullRequest, error) {
	var request models.PullRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	if request.PullRequestID == "" {
		return nil, errors.New("pull request id is required")
	}

	return &request, nil
}

func CreatePullRequestReassign(r *http.Request) (*models.PullRequest, string, error) {
	var request ReassignRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, "", err
	}

	if request.PullRequestID == "" {
		return nil, "", errors.New("pull request id is required")
	}

	if request.OldUserID == "" {
		return nil, "", errors.New("old user id is required")
	}

	return &models.PullRequest{
		PullRequestID: request.PullRequestID,
	}, request.OldUserID, nil
}

type ReassignRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}
