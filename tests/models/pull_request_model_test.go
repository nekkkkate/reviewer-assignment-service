package models

import (
	"reviewer-assignment-service/internal/domain/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupPRTest(t *testing.T) (*models.User, *models.User, *models.User, *models.Team, *models.PullRequest) {
	t.Helper()

	author := models.NewUser("Author", "author@example.com", true, "Developers")
	author.SetId(1)

	reviewer1 := models.NewUser("Reviewer1", "reviewer1@example.com", true, "Developers")
	reviewer1.SetId(2)

	reviewer2 := models.NewUser("Reviewer2", "reviewer2@example.com", true, "Developers")
	reviewer2.SetId(3)

	reviewer3 := models.NewUser("Reviewer3", "reviewer3@example.com", true, "Developers")
	reviewer3.SetId(4)

	team := models.NewTeam("Developers")

	team.AddMember(models.NewTeamMember(author.ID, author.Name, author.IsActive))
	team.AddMember(models.NewTeamMember(reviewer1.ID, reviewer1.Name, reviewer1.IsActive))
	team.AddMember(models.NewTeamMember(reviewer2.ID, reviewer2.Name, reviewer2.IsActive))
	team.AddMember(models.NewTeamMember(reviewer3.ID, reviewer3.Name, reviewer3.IsActive))

	pr, err := models.NewPullRequest("Test PR", author, team)
	require.NoError(t, err)

	return author, reviewer1, reviewer2, team, pr
}

func TestNewPullRequest_Success(t *testing.T) {
	author, _, _, team, _ := setupPRTest(t)

	pr, err := models.NewPullRequest("New Feature", author, team)
	require.NoError(t, err)

	assert.Equal(t, "New Feature", pr.Name)
	assert.Equal(t, models.StatusOpen, pr.Status)
	assert.Equal(t, author, pr.Author)
	assert.Empty(t, pr.Reviewers)
	assert.True(t, pr.CanModifyReviewers())
	assert.WithinDuration(t, time.Now(), pr.CreatedAt, time.Second)
}

func TestNewPullRequest_AuthorNotInTeam(t *testing.T) {
	author := models.NewUser("Author", "author@example.com", true, "OtherTeam")
	author.SetId(1)

	team := models.NewTeam("Developers") // Empty team

	_, err := models.NewPullRequest("Test PR", author, team)
	assert.Equal(t, models.ErrAuthorNotInTeam, err)
}

func TestPullRequest_AddReviewer_Success(t *testing.T) {
	_, reviewer1, reviewer2, _, pr := setupPRTest(t)

	err := pr.AddReviewer(reviewer1)
	require.NoError(t, err)
	assert.Len(t, pr.Reviewers, 1)
	assert.Equal(t, reviewer1, pr.Reviewers[0])

	err = pr.AddReviewer(reviewer2)
	require.NoError(t, err)
	assert.Len(t, pr.Reviewers, 2)
}

func TestPullRequest_AddReviewer_AlreadyAssigned(t *testing.T) {
	_, reviewer1, _, _, pr := setupPRTest(t)

	err := pr.AddReviewer(reviewer1)
	require.NoError(t, err)

	err = pr.AddReviewer(reviewer1)
	assert.Equal(t, models.ErrReviewerAlreadyAssigned, err)
	assert.Len(t, pr.Reviewers, 1)
}

func TestPullRequest_AddReviewer_TooManyReviewers(t *testing.T) {
	_, reviewer1, reviewer2, _, pr := setupPRTest(t)

	reviewer3 := models.NewUser("Reviewer3", "reviewer3@example.com", true, "Developers")
	reviewer3.SetId(4)

	err := pr.AddReviewer(reviewer1)
	require.NoError(t, err)

	err = pr.AddReviewer(reviewer2)
	require.NoError(t, err)

	err = pr.AddReviewer(reviewer3)
	assert.Equal(t, models.ErrTooManyReviewers, err)
	assert.Len(t, pr.Reviewers, 2)
}

func TestPullRequest_AddReviewer_MergedPR(t *testing.T) {
	_, reviewer1, _, _, pr := setupPRTest(t)

	pr.SetStatusMerged()

	err := pr.AddReviewer(reviewer1)
	assert.Equal(t, models.ErrPRAlreadyMerged, err)
	assert.Empty(t, pr.Reviewers)
}

func TestPullRequest_RemoveReviewer_Success(t *testing.T) {
	_, reviewer1, reviewer2, _, pr := setupPRTest(t)

	pr.AddReviewer(reviewer1)
	pr.AddReviewer(reviewer2)

	err := pr.RemoveReviewer(reviewer1.ID)
	require.NoError(t, err)
	assert.Len(t, pr.Reviewers, 1)
	assert.Equal(t, reviewer2, pr.Reviewers[0])
}

func TestPullRequest_RemoveReviewer_NotFound(t *testing.T) {
	_, _, _, _, pr := setupPRTest(t)

	err := pr.RemoveReviewer(999)
	assert.Equal(t, models.ErrReviewerNotFound, err)
}

func TestPullRequest_RemoveReviewer_MergedPR(t *testing.T) {
	_, reviewer1, _, _, pr := setupPRTest(t)

	pr.AddReviewer(reviewer1)
	pr.SetStatusMerged()

	err := pr.RemoveReviewer(reviewer1.ID)
	assert.Equal(t, models.ErrPRAlreadyMerged, err)
	assert.Len(t, pr.Reviewers, 1)
}

func TestPullRequest_ReplaceReviewer_Success(t *testing.T) {
	_, reviewer1, reviewer2, _, pr := setupPRTest(t)

	reviewer3 := models.NewUser("Reviewer3", "reviewer3@example.com", true, "Developers")
	reviewer3.SetId(4)

	pr.AddReviewer(reviewer1)
	pr.AddReviewer(reviewer2)

	err := pr.ReplaceReviewer(reviewer1.ID, reviewer3)
	require.NoError(t, err)
	assert.Len(t, pr.Reviewers, 2)

	found := false
	for _, reviewer := range pr.Reviewers {
		if reviewer.ID == reviewer3.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "New reviewer should be in the list")
}

func TestPullRequest_ReplaceReviewer_OldReviewerNotFound(t *testing.T) {
	_, reviewer1, _, _, pr := setupPRTest(t)

	err := pr.ReplaceReviewer(999, reviewer1)
	assert.Equal(t, models.ErrReviewerNotFound, err)
}

func TestPullRequest_Setters(t *testing.T) {
	_, _, _, _, pr := setupPRTest(t)

	pr.SetId(1)
	assert.Equal(t, 1, pr.ID)

	assert.True(t, pr.CanModifyReviewers())

	pr.SetStatusMerged()
	assert.Equal(t, models.StatusMerged, pr.Status)
	assert.False(t, pr.CanModifyReviewers())

	mergedTime := time.Now()
	pr.SetMergedAt(mergedTime)
	assert.Equal(t, mergedTime, pr.MergedAt)
}

func TestPullRequest_EdgeCases(t *testing.T) {
	author, _, _, team, _ := setupPRTest(t)

	pr, err := models.NewPullRequest("", author, team)
	require.NoError(t, err)
	assert.Equal(t, "", pr.Name)
}
