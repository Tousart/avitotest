package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/tousart/avitotest/internal/codes"
	"github.com/tousart/avitotest/internal/models"
)

func WriteAPIError(w http.ResponseWriter, httpStatus int, code, message string) {
	errResp := models.ErrorResponse{
		Code:    code,
		Message: message,
	}

	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(errResp)
}

func GetStatusError(code string) int {
	switch code {
	case codes.ErrTeamExists:
		return http.StatusBadRequest
	case codes.ErrNotFound:
		return http.StatusNotFound
	case codes.ErrPRExists:
		return http.StatusConflict
	case codes.ErrPRMerged:
		return http.StatusConflict
	case codes.ErrNotAssigned:
		return http.StatusConflict
	case codes.ErrNoCandidate:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
