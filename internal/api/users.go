package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tousart/avitotest/internal/api/helpers"
	"github.com/tousart/avitotest/internal/api/types"
	"github.com/tousart/avitotest/internal/codes"
	"github.com/tousart/avitotest/internal/models"
	"github.com/tousart/avitotest/internal/usecase"
)

type Users struct {
	usersService usecase.UsersService
}

func CreateUsersAPI(usersService usecase.UsersService) *Users {
	return &Users{
		usersService: usersService,
	}
}

func (u *Users) usersSetIsActiveHandler(w http.ResponseWriter, r *http.Request) {
	user, err := types.CreateSetIsActiveRequest(r)
	if err != nil {
		helpers.WriteAPIError(w, http.StatusBadRequest, codes.ErrBadRequet, "bad request")
		return
	}

	errResp := u.usersService.SetIsActive(r.Context(), user)
	if errResp != nil {
		status := helpers.GetStatusError(errResp.Code)
		helpers.WriteAPIError(w, status, errResp.Code, errResp.Message)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (u *Users) usersGetReviewHandler(w http.ResponseWriter, r *http.Request) {
	pullRequests, userID, err := types.CreateGetReview(r)
	if err != nil {
		helpers.WriteAPIError(w, http.StatusBadRequest, codes.ErrBadRequet, "bad request")
		return
	}

	errResp := u.usersService.GetReview(r.Context(), pullRequests, userID)
	if errResp != nil {
		status := helpers.GetStatusError(errResp.Code)
		helpers.WriteAPIError(w, status, errResp.Code, errResp.Message)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*pullRequests)
}

func (u *Users) usersGetActivityHandler(w http.ResponseWriter, r *http.Request) {
	var usersActivity []models.UserActivity

	errResp := u.usersService.GetActivity(r.Context(), &usersActivity)
	if errResp != nil {
		status := helpers.GetStatusError(errResp.Code)
		helpers.WriteAPIError(w, status, errResp.Code, errResp.Message)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(usersActivity)
}

func (u *Users) WithUsersHandlers(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/setIsActive", u.usersSetIsActiveHandler)
		r.Get("/getReview", u.usersGetReviewHandler)
		r.Get("/getActivity", u.usersGetActivityHandler)
	})
}
