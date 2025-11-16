package validators

import (
	"reviewer-assignment-service/internal/app/transport/dtos"
	"reviewer-assignment-service/internal/domain/models"
	"strconv"
	"strings"
)

func ValidateCreatePullRequestRequest(req *dtos.CreatePullRequestRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return NewValidationError("pull request name is required")
	}

	if len(req.Name) < 2 || len(req.Name) > 200 {
		return NewValidationError("pull request name must be between 2 and 200 characters")
	}

	if req.AuthorID <= 0 {
		return NewValidationError("author_id must be positive")
	}

	if len(req.Reviewers) > 2 {
		return NewValidationError("cannot assign more than 2 reviewers")
	}

	return nil
}

func ValidateUpdatePullRequestRequest(req *dtos.UpdatePullRequestRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return NewValidationError("pull request name is required")
	}

	if len(req.Name) < 2 || len(req.Name) > 200 {
		return NewValidationError("pull request name must be between 2 and 200 characters")
	}

	validStatuses := map[models.PRStatus]bool{
		models.StatusOpen:   true,
		models.StatusMerged: true,
	}
	if !validStatuses[models.PRStatus(req.Status)] {
		return NewValidationError("invalid status. Must be 'OPEN' or 'MERGED'")
	}

	if len(req.Reviewers) > 2 {
		return NewValidationError("cannot assign more than 2 reviewers")
	}

	return nil
}

func ValidateReassignReviewersRequest(req *dtos.ReassignReviewersRequest) error {
	if req.OldReviewerID <= 0 {
		return NewValidationError("old_reviewer_id must be positive")
	}
	return nil
}

func ValidateAddReviewerRequest(req *dtos.AddReviewerRequest) error {
	if req.ReviewerID <= 0 {
		return NewValidationError("reviewer_id must be positive")
	}
	return nil
}

func ValidatePullRequestID(prIDStr string) (int, error) {
	if prIDStr == "" {
		return 0, NewValidationError("pull request id is required")
	}

	prID, err := strconv.Atoi(prIDStr)
	if err != nil {
		return 0, NewValidationError("pull request id must be a valid number")
	}

	if prID <= 0 {
		return 0, NewValidationError("pull request id must be positive")
	}

	return prID, nil
}

func ValidateAuthorID(authorIDStr string) (int, error) {
	if authorIDStr == "" {
		return 0, NewValidationError("author id is required")
	}

	authorID, err := strconv.Atoi(authorIDStr)
	if err != nil {
		return 0, NewValidationError("author id must be a valid number")
	}

	if authorID <= 0 {
		return 0, NewValidationError("author id must be positive")
	}

	return authorID, nil
}

func ValidateReviewerID(reviewerIDStr string) (int, error) {
	if reviewerIDStr == "" {
		return 0, NewValidationError("reviewer id is required")
	}

	reviewerID, err := strconv.Atoi(reviewerIDStr)
	if err != nil {
		return 0, NewValidationError("reviewer id must be a valid number")
	}

	if reviewerID <= 0 {
		return 0, NewValidationError("reviewer id must be positive")
	}

	return reviewerID, nil
}

func ValidateStatus(status string) error {
	validStatuses := map[string]bool{
		"OPEN":   true,
		"MERGED": true,
	}
	if !validStatuses[status] {
		return NewValidationError("invalid status. Must be 'OPEN' or 'MERGED'")
	}
	return nil
}
