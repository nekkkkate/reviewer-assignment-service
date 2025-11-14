package persistence

import (
	"reviewer-assignment-service/internal/infrustructure/persistence/in-memory"
	"testing"

	"reviewer-assignment-service/internal/domain/models"
	"reviewer-assignment-service/internal/domain/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestData() (*models.User, *models.User, *models.User, *models.Team) {
	author := models.NewUser("Author", "author@example.com", true)
	author.SetId(1)

	reviewer1 := models.NewUser("Reviewer 1", "reviewer1@example.com", true)
	reviewer1.SetId(2)

	reviewer2 := models.NewUser("Reviewer 2", "reviewer2@example.com", true)
	reviewer2.SetId(3)

	reviewer3 := models.NewUser("Reviewer 3", "reviewer3@example.com", true)
	reviewer3.SetId(4)

	team := models.NewTeam("Development Team")
	team.SetId(1)
	team.AddUser(author)
	team.AddUser(reviewer1)
	team.AddUser(reviewer2)
	team.AddUser(reviewer3)

	return author, reviewer1, reviewer2, team
}

func createTestPR(id int, name string, author *models.User, team *models.Team) *models.PullRequest {
	pr, err := models.NewPullRequest(name, author, team)
	if err != nil {
		panic(err)
	}
	pr.SetId(id)
	return pr
}

func TestPullRequestRepository_Add_Success(t *testing.T) {
	repo := in_memory.NewPullRequestRepository()
	author, _, _, team := setupTestData()
	pr := createTestPR(1, "Test PR", author, team)

	err := repo.Add(pr)

	require.NoError(t, err)

	retrieved, err := repo.GetByID(1)
	require.NoError(t, err)
	assert.Equal(t, pr, retrieved)
}

func TestPullRequestRepository_Add_Duplicate(t *testing.T) {
	repo := in_memory.NewPullRequestRepository()
	author, _, _, team := setupTestData()
	pr := createTestPR(1, "Test PR", author, team)

	err := repo.Add(pr)
	require.NoError(t, err)

	err = repo.Add(pr)

	assert.Equal(t, repositories.ErrPullRequestAlreadyExists, err)
}

func TestPullRequestRepository_Add_Multiple(t *testing.T) {
	repo := in_memory.NewPullRequestRepository()
	author, _, _, team := setupTestData()

	pr1 := createTestPR(1, "PR 1", author, team)
	pr2 := createTestPR(2, "PR 2", author, team)
	pr3 := createTestPR(3, "PR 3", author, team)

	err := repo.Add(pr1)
	require.NoError(t, err)

	err = repo.Add(pr2)
	require.NoError(t, err)

	err = repo.Add(pr3)
	require.NoError(t, err)

	allPRs, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, allPRs, 3)
}

func TestPullRequestRepository_GetByID_Success(t *testing.T) {
	repo := in_memory.NewPullRequestRepository()
	author, _, _, team := setupTestData()
	pr := createTestPR(1, "Test PR", author, team)

	repo.Add(pr)

	retrieved, err := repo.GetByID(1)

	require.NoError(t, err)
	assert.Equal(t, pr, retrieved)
	assert.Equal(t, 1, retrieved.ID)
	assert.Equal(t, "Test PR", retrieved.Name)
	assert.Equal(t, models.StatusOpen, retrieved.Status)
}

func TestPullRequestRepository_GetByID_NotFound(t *testing.T) {
	repo := in_memory.NewPullRequestRepository()

	_, err := repo.GetByID(999)

	assert.Equal(t, repositories.ErrPullRequestNotFoundInPersistence, err)
}

func TestPullRequestRepository_GetByID_InvalidID(t *testing.T) {
	repo := in_memory.NewPullRequestRepository()

	_, err := repo.GetByID(0)
	assert.Equal(t, repositories.ErrPullRequestNotFoundInPersistence, err)

	_, err = repo.GetByID(-1)
	assert.Equal(t, repositories.ErrPullRequestNotFoundInPersistence, err)
}

func TestPullRequestRepository_GetAll(t *testing.T) {
	repo := in_memory.NewPullRequestRepository()
	author, _, _, team := setupTestData()

	pr1 := createTestPR(1, "PR 1", author, team)
	pr2 := createTestPR(2, "PR 2", author, team)
	pr3 := createTestPR(3, "PR 3", author, team)

	repo.Add(pr1)
	repo.Add(pr2)
	repo.Add(pr3)

	allPRs, err := repo.GetAll()

	require.NoError(t, err)
	assert.Len(t, allPRs, 3)

	prMap := make(map[int]*models.PullRequest)
	for _, pr := range allPRs {
		prMap[pr.ID] = pr
	}

	assert.Equal(t, pr1, prMap[1])
	assert.Equal(t, pr2, prMap[2])
	assert.Equal(t, pr3, prMap[3])
}

func TestPullRequestRepository_GetAll_Empty(t *testing.T) {
	repo := in_memory.NewPullRequestRepository()

	allPRs, err := repo.GetAll()

	require.NoError(t, err)
	assert.Empty(t, allPRs)
}

func TestPullRequestRepository_GetByStatus(t *testing.T) {
	repo := in_memory.NewPullRequestRepository()
	author, _, _, team := setupTestData()

	openPR1 := createTestPR(1, "Open PR 1", author, team)
	openPR2 := createTestPR(2, "Open PR 2", author, team)
	mergedPR := createTestPR(3, "Merged PR", author, team)
	mergedPR.SetStatusMerged()

	repo.Add(openPR1)
	repo.Add(openPR2)
	repo.Add(mergedPR)

	openPRs, err := repo.GetByStatus(models.StatusOpen)
	require.NoError(t, err)
	assert.Len(t, openPRs, 2)

	for _, pr := range openPRs {
		assert.Equal(t, models.StatusOpen, pr.Status)
	}

	mergedPRs, err := repo.GetByStatus(models.StatusMerged)
	require.NoError(t, err)
	assert.Len(t, mergedPRs, 1)
	assert.Equal(t, mergedPR, mergedPRs[0])
}

func TestPullRequestRepository_GetByStatus_NoMatches(t *testing.T) {
	repo := in_memory.NewPullRequestRepository()
	author, _, _, team := setupTestData()

	openPR := createTestPR(1, "Open PR", author, team)
	repo.Add(openPR)

	mergedPRs, err := repo.GetByStatus(models.StatusMerged)
	require.NoError(t, err)
	assert.Empty(t, mergedPRs)
}

func TestPullRequestRepository_GetByAuthorID(t *testing.T) {
	repo := in_memory.NewPullRequestRepository()
	author, _, _, team := setupTestData()

	otherAuthor := models.NewUser("Other Author", "other@example.com", true)
	otherAuthor.SetId(99)
	team.AddUser(otherAuthor)

	pr1 := createTestPR(1, "PR by Author", author, team)
	pr2 := createTestPR(2, "Another PR by Author", author, team)
	pr3 := createTestPR(3, "PR by Other Author", otherAuthor, team)

	repo.Add(pr1)
	repo.Add(pr2)
	repo.Add(pr3)

	authorPRs, err := repo.GetByAuthorID(author.ID)

	require.NoError(t, err)
	assert.Len(t, authorPRs, 2)

	for _, pr := range authorPRs {
		assert.Equal(t, author.ID, pr.Author.ID)
	}

	otherAuthorPRs, err := repo.GetByAuthorID(otherAuthor.ID)
	require.NoError(t, err)
	assert.Len(t, otherAuthorPRs, 1)
	assert.Equal(t, pr3, otherAuthorPRs[0])
}

func TestPullRequestRepository_GetByAuthorID_NotFound(t *testing.T) {
	repo := in_memory.NewPullRequestRepository()
	author, _, _, team := setupTestData()

	pr := createTestPR(1, "Test PR", author, team)
	repo.Add(pr)

	prs, err := repo.GetByAuthorID(999)

	require.NoError(t, err)
	assert.Empty(t, prs)
}

func TestPullRequestRepository_GetByReviewerID(t *testing.T) {
	repo := in_memory.NewPullRequestRepository()
	author, reviewer1, reviewer2, team := setupTestData()

	pr1 := createTestPR(1, "PR with Reviewer 1", author, team)
	pr1.AddReviewer(reviewer1)

	pr2 := createTestPR(2, "PR with Reviewer 2", author, team)
	pr2.AddReviewer(reviewer2)

	pr3 := createTestPR(3, "PR with Both Reviewers", author, team)
	pr3.AddReviewer(reviewer1)
	pr3.AddReviewer(reviewer2)

	repo.Add(pr1)
	repo.Add(pr2)
	repo.Add(pr3)

	reviewer1PRs, err := repo.GetByReviewerID(reviewer1.ID)
	require.NoError(t, err)
	assert.Len(t, reviewer1PRs, 2)

	for _, pr := range reviewer1PRs {
		found := false
		for _, reviewer := range pr.Reviewers {
			if reviewer.ID == reviewer1.ID {
				found = true
				break
			}
		}
		assert.True(t, found, "PR should contain reviewer1")
	}

	reviewer2PRs, err := repo.GetByReviewerID(reviewer2.ID)
	require.NoError(t, err)
	assert.Len(t, reviewer2PRs, 2)
}

func TestPullRequestRepository_GetByReviewerID_NoReviewers(t *testing.T) {
	repo := in_memory.NewPullRequestRepository()
	author, _, _, team := setupTestData()

	pr := createTestPR(1, "PR without reviewers", author, team)
	repo.Add(pr)

	prs, err := repo.GetByReviewerID(2)

	require.NoError(t, err)
	assert.Empty(t, prs)
}

func TestPullRequestRepository_GetByReviewerID_NotFound(t *testing.T) {
	repo := in_memory.NewPullRequestRepository()
	author, reviewer1, _, team := setupTestData()

	pr := createTestPR(1, "PR with reviewer", author, team)
	pr.AddReviewer(reviewer1)
	repo.Add(pr)

	prs, err := repo.GetByReviewerID(999)

	require.NoError(t, err)
	assert.Empty(t, prs)
}

func TestPullRequestRepository_ComplexScenario(t *testing.T) {
	repo := in_memory.NewPullRequestRepository()
	author, reviewer1, reviewer2, team := setupTestData()

	openPR1 := createTestPR(1, "Open PR 1", author, team)

	openPR2 := createTestPR(2, "Open PR 2", author, team)
	openPR2.AddReviewer(reviewer1)

	openPR3 := createTestPR(3, "Open PR 3", author, team)
	openPR3.AddReviewer(reviewer1)
	openPR3.AddReviewer(reviewer2)

	mergedPR := createTestPR(4, "Merged PR", author, team)
	mergedPR.SetStatusMerged()
	mergedPR.AddReviewer(reviewer2)

	repo.Add(openPR1)
	repo.Add(openPR2)
	repo.Add(openPR3)
	repo.Add(mergedPR)

	t.Run("GetAll returns all PRs", func(t *testing.T) {
		allPRs, err := repo.GetAll()
		require.NoError(t, err)
		assert.Len(t, allPRs, 4)
	})

	t.Run("GetByStatus filters correctly", func(t *testing.T) {
		openPRs, err := repo.GetByStatus(models.StatusOpen)
		require.NoError(t, err)
		assert.Len(t, openPRs, 3)

		mergedPRs, err := repo.GetByStatus(models.StatusMerged)
		require.NoError(t, err)
		assert.Len(t, mergedPRs, 1)
	})

	t.Run("GetByAuthorID finds author's PRs", func(t *testing.T) {
		authorPRs, err := repo.GetByAuthorID(author.ID)
		require.NoError(t, err)
		assert.Len(t, authorPRs, 4)
	})

}

func TestPullRequestRepository_ConcurrentAccess(t *testing.T) {
	repo := in_memory.NewPullRequestRepository()
	author, _, _, team := setupTestData()

	pr := createTestPR(1, "Concurrent Test PR", author, team)

	err := repo.Add(pr)
	require.NoError(t, err)

	retrieved, err := repo.GetByID(1)
	require.NoError(t, err)
	assert.Equal(t, pr, retrieved)

	openPRs, err := repo.GetByStatus(models.StatusOpen)
	require.NoError(t, err)
	assert.Len(t, openPRs, 1)

	authorPRs, err := repo.GetByAuthorID(author.ID)
	require.NoError(t, err)
	assert.Len(t, authorPRs, 1)
}

func TestPullRequestRepository_EdgeCases(t *testing.T) {
	repo := in_memory.NewPullRequestRepository()

	t.Run("Empty repository operations", func(t *testing.T) {
		_, err := repo.GetByID(1)
		assert.Equal(t, repositories.ErrPullRequestNotFoundInPersistence, err)

		all, err := repo.GetAll()
		require.NoError(t, err)
		assert.Empty(t, all)

		byStatus, err := repo.GetByStatus(models.StatusOpen)
		require.NoError(t, err)
		assert.Empty(t, byStatus)

		byAuthor, err := repo.GetByAuthorID(1)
		require.NoError(t, err)
		assert.Empty(t, byAuthor)

		byReviewer, err := repo.GetByReviewerID(1)
		require.NoError(t, err)
		assert.Empty(t, byReviewer)
	})
}
