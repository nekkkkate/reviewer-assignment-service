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
	switch {
	case errors.Is(err, repositories.ErrUserAlreadyExists):
		SendError(w, "USER_ALREADY_EXISTS", "User already exists", http.StatusConflict)
	case errors.Is(err, repositories.ErrUserNotFoundInPersistence):
		SendError(w, "USER_NOT_FOUND", "User not found", http.StatusNotFound)
	case errors.Is(err, repositories.ErrUserWithThatEmailNotFound):
		SendError(w, "USER_NOT_FOUND", "User with this email not found", http.StatusNotFound)

	case errors.Is(err, repositories.ErrTeamNotFoundInPersistence):
		SendError(w, "TEAM_NOT_FOUND", "Team not found", http.StatusNotFound)
	case errors.Is(err, repositories.ErrTeamAlreadyExists):
		SendError(w, "TEAM_ALREADY_EXISTS", "Team already exists", http.StatusConflict)

	case errors.Is(err, models.ErrMemberAlreadyInTeam):
		SendError(w, "MEMBER_ALREADY_IN_TEAM", "User is already a member of this team", http.StatusConflict)
	case errors.Is(err, models.ErrMemberNotInTeam):
		SendError(w, "MEMBER_NOT_IN_TEAM", "User is not a member of this team", http.StatusNotFound)

	case errors.Is(err, models.ErrAuthorNotInTeam):
		SendError(w, "AUTHOR_NOT_IN_TEAM", "Author not in team", http.StatusBadRequest)
	case errors.Is(err, models.ErrPRAlreadyMerged):
		SendError(w, "PR_ALREADY_MERGED", "Cannot reassign on merged PR", http.StatusConflict)
	case errors.Is(err, models.ErrReviewerNotFound):
		SendError(w, "REVIEWER_NOT_FOUND", "No active replacement candidate in team", http.StatusNotFound)
	case errors.Is(err, models.ErrReviewerAlreadyAssigned):
		SendError(w, "REVIEWER_ALREADY_ASSIGNED", "Reviewer already assigned to this PR", http.StatusConflict)
	case errors.Is(err, models.ErrTooManyReviewers):
		SendError(w, "TOO_MANY_REVIEWERS", "Too many reviewers assigned", http.StatusBadRequest)
	case errors.Is(err, repositories.ErrPullRequestNotFoundInPersistence):
		SendError(w, "PR_NOT_FOUND", "Pull request not found", http.StatusNotFound)
	case errors.Is(err, repositories.ErrPullRequestAlreadyExists):
		SendError(w, "PR_ALREADY_EXISTS", "Pull request already exists", http.StatusConflict)

	case isValidationError(err):
		SendError(w, "VALIDATION_ERROR", err.Error(), http.StatusBadRequest)

	default:
		SendError(w, "INTERNAL_ERROR", "Internal server error", http.StatusInternalServerError)
	}
}

func isValidationError(err error) bool {
	var validationErr *validators.ValidationError
	return errors.As(err, &validationErr)
}

func SendValidationError(w http.ResponseWriter, message string) {
	SendError(w, "VALIDATION_ERROR", message, http.StatusBadRequest)
}

func SendNotFound(w http.ResponseWriter, resource string) {
	SendError(w, "NOT_FOUND", resource+" not found", http.StatusNotFound)
}

func SendBadRequest(w http.ResponseWriter, message string) {
	SendError(w, "BAD_REQUEST", message, http.StatusBadRequest)
}

func SendInternalError(w http.ResponseWriter) {
	SendError(w, "INTERNAL_ERROR", "Internal server error", http.StatusInternalServerError)
}
