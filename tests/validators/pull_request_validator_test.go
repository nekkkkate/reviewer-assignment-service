package validators

import (
	"testing"

	"reviewer-assignment-service/internal/app/transport/dtos"
	"reviewer-assignment-service/internal/app/validators"

	"github.com/stretchr/testify/assert"
)

func TestValidateCreatePullRequestRequest_InvalidCases(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		req := &dtos.CreatePullRequestRequest{
			Name:      "",
			AuthorID:  1,
			Reviewers: []int{},
		}

		err := validators.ValidateCreatePullRequestRequest(req)
		if assert.Error(t, err) {
			assert.Equal(t, "pull request name is required", err.Error())
		}
	})

	t.Run("name too short", func(t *testing.T) {
		req := &dtos.CreatePullRequestRequest{
			Name:      "a",
			AuthorID:  1,
			Reviewers: []int{},
		}

		err := validators.ValidateCreatePullRequestRequest(req)
		if assert.Error(t, err) {
			assert.Equal(t, "pull request name must be between 2 and 200 characters", err.Error())
		}
	})

	t.Run("author id not positive", func(t *testing.T) {
		req := &dtos.CreatePullRequestRequest{
			Name:      "Valid Name",
			AuthorID:  0,
			Reviewers: []int{},
		}

		err := validators.ValidateCreatePullRequestRequest(req)
		if assert.Error(t, err) {
			assert.Equal(t, "author_id must be positive", err.Error())
		}
	})

	t.Run("too many reviewers", func(t *testing.T) {
		req := &dtos.CreatePullRequestRequest{
			Name:      "Valid Name",
			AuthorID:  1,
			Reviewers: []int{1, 2, 3},
		}

		err := validators.ValidateCreatePullRequestRequest(req)
		if assert.Error(t, err) {
			assert.Equal(t, "cannot assign more than 2 reviewers", err.Error())
		}
	})
}

func TestValidateCreatePullRequestRequest_ValidCase(t *testing.T) {
	req := &dtos.CreatePullRequestRequest{
		Name:      "Valid PR",
		AuthorID:  1,
		Reviewers: []int{2},
	}

	err := validators.ValidateCreatePullRequestRequest(req)
	assert.NoError(t, err)
}
