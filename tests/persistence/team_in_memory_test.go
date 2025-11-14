package persistence

import (
	"reviewer-assignment-service/internal/domain/models"
	"reviewer-assignment-service/internal/domain/repositories"
	"reviewer-assignment-service/internal/infrustructure/persistence/in-memory"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTeamRepository(t *testing.T) *in_memory.TeamRepository {
	t.Helper()
	return in_memory.NewTeamRepository()
}

func createTestTeam(id int, name string) *models.Team {
	team := models.NewTeam(name)
	team.SetId(id)
	return team
}

func TestTeamRepository_Add_Success(t *testing.T) {
	repo := setupTeamRepository(t)
	team := createTestTeam(1, "Development Team")

	err := repo.Add(team)
	require.NoError(t, err)

	retrieved, err := repo.GetByID(1)
	require.NoError(t, err)
	assert.Equal(t, team, retrieved)
}

func TestTeamRepository_Add_Duplicate(t *testing.T) {
	repo := setupTeamRepository(t)
	team := createTestTeam(1, "Development Team")

	err := repo.Add(team)
	require.NoError(t, err)

	err = repo.Add(team)
	assert.Equal(t, repositories.ErrTeamAlreadyExists, err)
}

func TestTeamRepository_GetByID_Success(t *testing.T) {
	repo := setupTeamRepository(t)
	team := createTestTeam(1, "Development Team")

	repo.Add(team)

	retrieved, err := repo.GetByID(1)
	require.NoError(t, err)
	assert.Equal(t, team, retrieved)
}

func TestTeamRepository_GetByID_NotFound(t *testing.T) {
	repo := setupTeamRepository(t)

	_, err := repo.GetByID(999)
	assert.Equal(t, repositories.ErrTeamNotFoundInPersistence, err)
}

func TestTeamRepository_GetByName_Success(t *testing.T) {
	repo := setupTeamRepository(t)
	team := createTestTeam(1, "Development Team")

	repo.Add(team)

	retrieved, err := repo.GetByName("Development Team")
	require.NoError(t, err)
	assert.Equal(t, team, retrieved)
}

func TestTeamRepository_GetByName_NotFound(t *testing.T) {
	repo := setupTeamRepository(t)

	_, err := repo.GetByName("Nonexistent Team")
	assert.Equal(t, repositories.ErrTeamNotFoundInPersistence, err)
}

func TestTeamRepository_GetByName_CaseSensitive(t *testing.T) {
	repo := setupTeamRepository(t)
	team := createTestTeam(1, "Development Team")

	repo.Add(team)

	_, err := repo.GetByName("development team")
	assert.Equal(t, repositories.ErrTeamNotFoundInPersistence, err)

	retrieved, err := repo.GetByName("Development Team")
	require.NoError(t, err)
	assert.Equal(t, team, retrieved)
}

func TestTeamRepository_GetAll(t *testing.T) {
	repo := setupTeamRepository(t)

	team1 := createTestTeam(1, "Team 1")
	team2 := createTestTeam(2, "Team 2")
	team3 := createTestTeam(3, "Team 3")

	repo.Add(team1)
	repo.Add(team2)
	repo.Add(team3)

	teams, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, teams, 3)

	teamMap := make(map[int]*models.Team)
	for _, team := range teams {
		teamMap[team.ID] = team
	}

	assert.Equal(t, team1, teamMap[1])
	assert.Equal(t, team2, teamMap[2])
	assert.Equal(t, team3, teamMap[3])
}

func TestTeamRepository_GetAll_Empty(t *testing.T) {
	repo := setupTeamRepository(t)

	teams, err := repo.GetAll()
	require.NoError(t, err)
	assert.Empty(t, teams)
}

func TestTeamRepository_Update_Success(t *testing.T) {
	repo := setupTeamRepository(t)
	team := createTestTeam(1, "Old Team Name")

	repo.Add(team)

	team.UpdateName("New Team Name")
	user := models.NewUser("Test User", "test@example.com", true)
	user.SetId(1)
	team.AddUser(user)

	err := repo.Update(team)
	require.NoError(t, err)

	updated, err := repo.GetByID(1)
	require.NoError(t, err)
	assert.Equal(t, "New Team Name", updated.Name)
	assert.True(t, updated.IsUserInTeam(user))
}

func TestTeamRepository_Update_NotFound(t *testing.T) {
	repo := setupTeamRepository(t)
	team := createTestTeam(999, "Nonexistent Team")

	err := repo.Update(team)
	assert.Equal(t, repositories.ErrTeamNotFoundInPersistence, err)
}

func TestTeamRepository_ComplexScenario(t *testing.T) {
	repo := setupTeamRepository(t)

	devTeam := createTestTeam(1, "Development")
	qaTeam := createTestTeam(2, "QA")
	opsTeam := createTestTeam(3, "Operations")

	repo.Add(devTeam)
	repo.Add(qaTeam)
	repo.Add(opsTeam)

	teams, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, teams, 3)

	foundTeam, err := repo.GetByName("QA")
	require.NoError(t, err)
	assert.Equal(t, "QA", foundTeam.Name)

	qaTeam.UpdateName("Quality Assurance")
	err = repo.Update(qaTeam)
	require.NoError(t, err)

	updated, err := repo.GetByID(2)
	require.NoError(t, err)
	assert.Equal(t, "Quality Assurance", updated.Name)

	_, err = repo.GetByName("QA")
	assert.Equal(t, repositories.ErrTeamNotFoundInPersistence, err)
}
