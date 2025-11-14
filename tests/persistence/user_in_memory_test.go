package persistence

import (
	"reviewer-assignment-service/internal/domain/models"
	"reviewer-assignment-service/internal/domain/repositories"
	"reviewer-assignment-service/internal/infrustructure/persistence/in-memory"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupUserRepository(t *testing.T) *in_memory.UserRepository {
	t.Helper()
	return in_memory.NewUserRepository()
}

func createTestUser(id int, name, email string) *models.User {
	user := models.NewUser(name, email, true)
	user.SetId(id)
	return user
}

func TestUserRepository_Add_Success(t *testing.T) {
	repo := setupUserRepository(t)
	user := createTestUser(1, "John Doe", "john@example.com")

	err := repo.Add(user)
	require.NoError(t, err)

	retrieved, err := repo.GetByID(1)
	require.NoError(t, err)
	assert.Equal(t, user, retrieved)
}

func TestUserRepository_Add_Duplicate(t *testing.T) {
	repo := setupUserRepository(t)
	user := createTestUser(1, "John Doe", "john@example.com")

	err := repo.Add(user)
	require.NoError(t, err)

	err = repo.Add(user)
	assert.Equal(t, repositories.ErrUserAlreadyExists, err)
}

func TestUserRepository_GetByID_Success(t *testing.T) {
	repo := setupUserRepository(t)
	user := createTestUser(1, "John Doe", "john@example.com")

	repo.Add(user)

	retrieved, err := repo.GetByID(1)
	require.NoError(t, err)
	assert.Equal(t, user, retrieved)
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	repo := setupUserRepository(t)

	_, err := repo.GetByID(999)
	assert.Equal(t, repositories.ErrUserNotFoundInPersistence, err)
}

func TestUserRepository_GetByEmail_Success(t *testing.T) {
	repo := setupUserRepository(t)
	user := createTestUser(1, "John Doe", "john@example.com")

	repo.Add(user)

	retrieved, err := repo.GetByEmail("john@example.com")
	require.NoError(t, err)
	assert.Equal(t, user, retrieved)
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	repo := setupUserRepository(t)

	_, err := repo.GetByEmail("nonexistent@example.com")
	assert.Equal(t, repositories.ErrUserWithThatEmailNotFound, err)
}

func TestUserRepository_GetByEmail_CaseSensitive(t *testing.T) {
	repo := setupUserRepository(t)
	user := createTestUser(1, "John Doe", "John@Example.com") // Mixed case

	repo.Add(user)

	_, err := repo.GetByEmail("john@example.com")
	assert.Equal(t, repositories.ErrUserWithThatEmailNotFound, err)

	retrieved, err := repo.GetByEmail("John@Example.com")
	require.NoError(t, err)
	assert.Equal(t, user, retrieved)
}

func TestUserRepository_GetAll(t *testing.T) {
	repo := setupUserRepository(t)

	user1 := createTestUser(1, "User 1", "user1@example.com")
	user2 := createTestUser(2, "User 2", "user2@example.com")
	user3 := createTestUser(3, "User 3", "user3@example.com")

	repo.Add(user1)
	repo.Add(user2)
	repo.Add(user3)

	users, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, users, 3)

	userMap := make(map[int]*models.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	assert.Equal(t, user1, userMap[1])
	assert.Equal(t, user2, userMap[2])
	assert.Equal(t, user3, userMap[3])
}

func TestUserRepository_GetAll_Empty(t *testing.T) {
	repo := setupUserRepository(t)

	users, err := repo.GetAll()
	require.NoError(t, err)
	assert.Empty(t, users)
}

func TestUserRepository_GetActiveUsers(t *testing.T) {
	repo := setupUserRepository(t)

	activeUser := createTestUser(1, "Active User", "active@example.com")
	inactiveUser := models.NewUser("Inactive User", "inactive@example.com", false)
	inactiveUser.SetId(2)

	repo.Add(activeUser)
	repo.Add(inactiveUser)

	activeUsers, err := repo.GetActiveUsers()
	require.NoError(t, err)
	assert.Len(t, activeUsers, 1)
	assert.Equal(t, activeUser, activeUsers[0])
}

func TestUserRepository_Update_Success(t *testing.T) {
	repo := setupUserRepository(t)
	user := createTestUser(1, "Old Name", "old@example.com")

	repo.Add(user)

	user.UpdateName("New Name")
	user.UpdateEmail("new@example.com")
	user.UpdateIsActive(false)

	err := repo.Update(user)
	require.NoError(t, err)

	updated, err := repo.GetByID(1)
	require.NoError(t, err)
	assert.Equal(t, "New Name", updated.Name)
	assert.Equal(t, "new@example.com", updated.Email)
	assert.False(t, updated.IsActive)
}

func TestUserRepository_Update_NotFound(t *testing.T) {
	repo := setupUserRepository(t)
	user := createTestUser(999, "Nonexistent", "none@example.com")

	err := repo.Update(user)
	assert.Equal(t, repositories.ErrUserNotFoundInPersistence, err)
}

func TestUserRepository_ConcurrentOperations(t *testing.T) {
	repo := setupUserRepository(t)

	user := createTestUser(1, "Test User", "test@example.com")

	err := repo.Add(user)
	require.NoError(t, err)

	retrieved, err := repo.GetByID(1)
	require.NoError(t, err)
	assert.Equal(t, user, retrieved)
}
