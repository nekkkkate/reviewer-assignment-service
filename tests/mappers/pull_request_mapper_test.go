package mappers

import (
	"testing"
	"time"

	"reviewer-assignment-service/internal/app/transport/mappers"
	"reviewer-assignment-service/internal/domain/models"

	"github.com/stretchr/testify/assert"
)

func TestToPullRequestResponse(t *testing.T) {
	author := &models.User{ID: 1, Name: "Author", Email: "author@example.com", TeamName: "backend", IsActive: true}
	reviewer1 := &models.User{ID: 2, Name: "Reviewer1", Email: "rev1@example.com", TeamName: "backend", IsActive: true}
	reviewer2 := &models.User{ID: 3, Name: "Reviewer2", Email: "rev2@example.com", TeamName: "backend", IsActive: true}

	createdAt := time.Now()
	mergedAt := createdAt.Add(time.Hour)

	pr := &models.PullRequest{
		ID:        10,
		Name:      "Test PR",
		Status:    models.StatusOpen,
		Author:    author,
		Reviewers: []*models.User{reviewer1, reviewer2},
		CreatedAt: createdAt,
		MergedAt:  mergedAt,
	}

	resp := mappers.ToPullRequestResponse(pr)

	assert.Equal(t, 10, resp.ID)
	assert.Equal(t, "Test PR", resp.Name)
	assert.Equal(t, string(models.StatusOpen), resp.Status)

	if assert.NotNil(t, resp.Author) {
		assert.Equal(t, "1", resp.Author.UserID)
		assert.Equal(t, "Author", resp.Author.Username)
		assert.Equal(t, "backend", resp.Author.TeamName)
		assert.True(t, resp.Author.IsActive)
	}

	if assert.Len(t, resp.Reviewers, 2) {
		assert.Equal(t, "2", resp.Reviewers[0].UserID)
		assert.Equal(t, "3", resp.Reviewers[1].UserID)
	}

	assert.Equal(t, createdAt, resp.CreatedAt)
	if assert.NotNil(t, resp.MergedAt) {
		assert.Equal(t, mergedAt, *resp.MergedAt)
	}
}
