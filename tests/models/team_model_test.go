package models

import (
	"reviewer-assignment-service/internal/domain/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTeamTest(t *testing.T) (*models.Team, *models.TeamMember, *models.TeamMember) {
	t.Helper()

	team := models.NewTeam("Developers")
	member1 := models.NewTeamMember(1, "John Doe", true)
	member2 := models.NewTeamMember(2, "Jane Doe", true)

	return team, member1, member2
}

func TestTeam_AddMember_Success(t *testing.T) {
	team, member1, _ := setupTeamTest(t)

	err := team.AddMember(member1)
	require.NoError(t, err)

	assert.True(t, team.IsMemberInTeam(member1.UserID))
	assert.Equal(t, 1, team.GetMemberCount())
	assert.False(t, team.IsEmpty())
}

func TestTeam_AddMember_AlreadyInTeam(t *testing.T) {
	team, member1, _ := setupTeamTest(t)

	err := team.AddMember(member1)
	require.NoError(t, err)

	err = team.AddMember(member1)
	assert.Equal(t, models.ErrMemberAlreadyInTeam, err)
	assert.Equal(t, 1, team.GetMemberCount())
}

func TestTeam_RemoveMember_Success(t *testing.T) {
	team, member1, _ := setupTeamTest(t)

	team.AddMember(member1)

	err := team.RemoveMember(member1)
	require.NoError(t, err)

	assert.False(t, team.IsMemberInTeam(member1.UserID))
	assert.True(t, team.IsEmpty())
	assert.Equal(t, 0, team.GetMemberCount())
}

func TestTeam_RemoveMember_NotInTeam(t *testing.T) {
	team, member1, _ := setupTeamTest(t)

	err := team.RemoveMember(member1)
	assert.Equal(t, models.ErrMemberNotInTeam, err)
}

func TestTeam_GetActiveMembers(t *testing.T) {
	team, member1, member2 := setupTeamTest(t)

	member3 := models.NewTeamMember(3, "Inactive Member", false)

	team.AddMember(member1)
	team.AddMember(member2)
	team.AddMember(member3)

	activeMembers := team.GetActiveMembers()
	assert.Len(t, activeMembers, 2)

	for _, member := range activeMembers {
		assert.True(t, member.IsActive)
	}
}

func TestTeam_Setters(t *testing.T) {
	team := models.NewTeam("Old Name")

	team.SetId(1)
	assert.Equal(t, 1, team.ID)

	team.UpdateName("New Name")
	assert.Equal(t, "New Name", team.Name)
}

func TestTeam_GetMembers(t *testing.T) {
	team, member1, member2 := setupTeamTest(t)

	team.AddMember(member1)
	team.AddMember(member2)

	members := team.GetMembers()
	assert.Len(t, members, 2)

	userIDs := make(map[int]bool)
	for _, member := range members {
		userIDs[member.UserID] = true
	}

	assert.True(t, userIDs[1])
	assert.True(t, userIDs[2])
}

func TestTeamMember_Setters(t *testing.T) {
	member := models.NewTeamMember(1, "Old Name", true)

	member.UpdateUsername("New Name")
	assert.Equal(t, "New Name", member.Username)

	member.UpdateIsActive(false)
	assert.False(t, member.IsActive)
}
