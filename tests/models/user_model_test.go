package models

import (
	"reviewer-assignment-service/internal/domain/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_CreationAndSetters(t *testing.T) {
	user := models.NewUser("John Doe", "john@example.com", true, "Developers")

	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@example.com", user.Email)
	assert.True(t, user.IsActive)
	assert.Equal(t, "Developers", user.TeamName)

	user.SetId(1)
	assert.Equal(t, 1, user.ID)

	user.UpdateName("Jane Doe")
	assert.Equal(t, "Jane Doe", user.Name)

	user.UpdateEmail("jane@example.com")
	assert.Equal(t, "jane@example.com", user.Email)

	user.UpdateIsActive(false)
	assert.False(t, user.IsActive)
}

func TestUser_EdgeCases(t *testing.T) {
	user := models.NewUser("", "", false, "")
	assert.Equal(t, "", user.Name)
	assert.Equal(t, "", user.Email)
	assert.False(t, user.IsActive)
	assert.Equal(t, "", user.TeamName)
}
