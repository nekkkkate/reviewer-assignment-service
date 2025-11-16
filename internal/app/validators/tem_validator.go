package validators

import (
	"reviewer-assignment-service/internal/app/transport/dtos"
	"strconv"
	"strings"
)

func ValidateCreateTeamRequest(req *dtos.CreateTeamRequest) error {
	if req.Name == "" {
		return NewValidationError("team name is required")
	}

	if len(req.Name) < 2 || len(req.Name) > 100 {
		return NewValidationError("team name must be between 2 and 100 characters")
	}

	for _, member := range req.Members {
		if member.UserID <= 0 {
			return NewValidationError("member %d: user_id must be positive")
		}
		if strings.TrimSpace(member.Username) == "" {
			return NewValidationError("member %d: username is required")
		}
		if len(member.Username) > 100 {
			return NewValidationError("member %d: username must be less than 100 characters")
		}
	}

	return nil
}

func ValidateUpdateTeamRequest(req *dtos.UpdateTeamRequest) error {
	if req.Name == "" {
		return NewValidationError("team name is required")
	}

	if len(req.Name) < 2 || len(req.Name) > 100 {
		return NewValidationError("team name must be between 2 and 100 characters")
	}

	for _, member := range req.Members {
		if member.UserID <= 0 {
			return NewValidationError("member %d: user_id must be positive")
		}
		if strings.TrimSpace(member.Username) == "" {
			return NewValidationError("member %d: username is required")
		}
	}

	return nil
}

func ValidateAddMemberRequest(req *dtos.AddMemberRequest) error {
	if req.UserID <= 0 {
		return NewValidationError("user_id must be positive")
	}
	return nil
}

func ValidateTeamID(teamIDStr string) (int, error) {
	if teamIDStr == "" {
		return 0, NewValidationError("team id is required")
	}

	teamID, err := strconv.Atoi(teamIDStr)
	if err != nil {
		return 0, NewValidationError("team id must be a valid number")
	}

	if teamID <= 0 {
		return 0, NewValidationError("team id must be positive")
	}

	return teamID, nil
}

func ValidateTeamName(name string) error {
	if name == "" {
		return NewValidationError("team name is required")
	}

	if len(name) < 2 || len(name) > 100 {
		return NewValidationError("team name must be between 2 and 100 characters")
	}

	return nil
}
