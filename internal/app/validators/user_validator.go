package validators

import (
	"reviewer-assignment-service/internal/app/transport/dtos"
	"strconv"
	"strings"
)

func ValidateSetUserActiveRequest(req *dtos.SetUserActiveRequest) error {
	if req.UserID == "" {
		return NewValidationError("user_id is required")
	}

	if _, err := strconv.Atoi(req.UserID); err != nil {
		return NewValidationError("user_id must be a valid number")
	}

	return nil
}

func ValidateUserID(userID string) error {
	if userID == "" {
		return NewValidationError("user_id is required")
	}

	if _, err := strconv.Atoi(userID); err != nil {
		return NewValidationError("user_id must be a valid number")
	}

	return nil
}

func ValidateCreateUserRequest(req *dtos.CreateUserRequest) error {
	if req.Username == "" {
		return NewValidationError("username is required")
	}

	if req.Email == "" {
		return NewValidationError("email is required")
	}

	if !strings.Contains(req.Email, "@") {
		return NewValidationError("email must be a valid email address")
	}

	if req.TeamName == "" {
		return NewValidationError("team_name is required")
	}

	return nil
}

func ValidateGetUserByEmailRequest(req *dtos.GetUserByEmailRequest) error {
	if req.Email == "" {
		return NewValidationError("email is required")
	}

	if !strings.Contains(req.Email, "@") {
		return NewValidationError("email must be a valid email address")
	}

	return nil
}

func ValidateDeactivateUserRequest(req *dtos.DeactivateUserRequest) error {
	if req.UserID == "" {
		return NewValidationError("user_id is required")
	}

	if _, err := strconv.Atoi(req.UserID); err != nil {
		return NewValidationError("user_id must be a valid number")
	}

	return nil
}

func ValidateEmail(email string) error {
	if email == "" {
		return NewValidationError("email is required")
	}

	if !strings.Contains(email, "@") {
		return NewValidationError("email must be a valid email address")
	}

	return nil
}

type ValidationError struct {
	Message string
}

func (v ValidationError) Error() string {
	return v.Message
}

func NewValidationError(message string) error {
	return &ValidationError{Message: message}
}
