package response_errors

import (
	"encoding/json"
	"errors"
	"net/http"
	"reviewer-assignment-service/internal/app/validators"
	"reviewer-assignment-service/internal/domain/models"
	"reviewer-assignment-service/internal/domain/repositories"
)

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func SendError(w http.ResponseWriter, code, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := ErrorResponse{}
	response.Error.Code = code
	response.Error.Message = message

	json.NewEncoder(w).Encode(response)
}

func HandleServiceError(w http.ResponseWriter, err error) {
	switch err {
	case models.ErrAuthorNotInTeam:
		SendError(w, "NOT_FOUND", "Author not in team", http.StatusNotFound)
	case models.ErrPRAlreadyMerged:
		SendError(w, "PR_MERGED", "Cannot reassign on merged PR", http.StatusConflict)
	case models.ErrReviewerNotFound:
		SendError(w, "NO_CANDIDATE", "No active replacement candidate in team", http.StatusConflict)
	case repositories.ErrUserAlreadyExists:
		SendError(w, "USER_EXISTS", "User already exists", http.StatusBadRequest)
	case repositories.ErrUserNotFoundInPersistence:
		SendError(w, "NOT_FOUND", "User not found", http.StatusNotFound)
	case repositories.ErrUserWithThatEmailNotFound:
		SendError(w, "NOT_FOUND", "User with this email not found", http.StatusNotFound)
	default:
		var validationError *validators.ValidationError
		if errors.As(err, &validationError) {
			SendError(w, "INVALID_REQUEST", err.Error(), http.StatusBadRequest)
		}
	}
}
