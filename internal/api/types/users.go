package types

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/tousart/avitotest/internal/models"
)

func CreateSetIsActiveRequest(r *http.Request) (*models.User, error) {
	var request *models.User

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	if request.UserID == "" {
		return nil, errors.New("user_id is required")
	}

	return request, nil
}

func CreateGetReview(r *http.Request) (*[]models.PullRequestShort, string, error) {
	userID := r.URL.Query().Get("user_id")

	if userID == "" {
		return nil, "", errors.New("user id is required")
	}

	pullRequests := make([]models.PullRequestShort, 0)

	return &pullRequests, userID, nil
}
