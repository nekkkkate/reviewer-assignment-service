package service

import (
	"reviewer-assignment-service/internal/domain/models"
	"reviewer-assignment-service/internal/domain/repositories"
	"reviewer-assignment-service/internal/domain/services/impl"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPullRequestRepository struct {
	mock.Mock
}

func (m *MockPullRequestRepository) Add(pr *models.PullRequest) error {
	args := m.Called(pr)
	return args.Error(0)
}

func (m *MockPullRequestRepository) GetByID(id int) (*models.PullRequest, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PullRequest), args.Error(1)
}

func (m *MockPullRequestRepository) GetAll() ([]*models.PullRequest, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.PullRequest), args.Error(1)
}

func (m *MockPullRequestRepository) GetByStatus(status models.PRStatus) ([]*models.PullRequest, error) {
	args := m.Called(status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.PullRequest), args.Error(1)
}

func (m *MockPullRequestRepository) GetByAuthorID(authorID int) ([]*models.PullRequest, error) {
	args := m.Called(authorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.PullRequest), args.Error(1)
}

func (m *MockPullRequestRepository) GetByReviewerID(reviewerID int) ([]*models.PullRequest, error) {
	args := m.Called(reviewerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.PullRequest), args.Error(1)
}

func (m *MockPullRequestRepository) Update(pr *models.PullRequest) error {
	args := m.Called(pr)
	return args.Error(0)
}

func (m *MockPullRequestRepository) FindPossibleReviewers(author *models.User) ([]*models.User, error) {
	args := m.Called(author)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func TestPullRequestService_Create(t *testing.T) {
	t.Run("successful PR creation", func(t *testing.T) {
		mockRepo := new(MockPullRequestRepository)
		prService := impl.NewPullRequestService(mockRepo)

		author := &models.User{
			ID:       1,
			Name:     "John Doe",
			Email:    "john@example.com",
			TeamName: "backend",
			IsActive: true,
		}

		pr := &models.PullRequest{
			Name:      "New feature",
			Status:    models.StatusOpen,
			Author:    author,
			Reviewers: []*models.User{},
			CreatedAt: time.Now(),
		}

		mockRepo.On("Add", pr).Return(nil)

		err := prService.Create(pr)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestPullRequestService_GetByID(t *testing.T) {
	t.Run("successful get by id", func(t *testing.T) {
		mockRepo := new(MockPullRequestRepository)
		prService := impl.NewPullRequestService(mockRepo)

		author := &models.User{
			ID:       1,
			Name:     "John Doe",
			Email:    "john@example.com",
			TeamName: "backend",
			IsActive: true,
		}

		expectedPR := &models.PullRequest{
			ID:        1,
			Name:      "New feature",
			Status:    models.StatusOpen,
			Author:    author,
			Reviewers: []*models.User{},
			CreatedAt: time.Now(),
		}

		mockRepo.On("GetByID", 1).Return(expectedPR, nil)

		pr, err := prService.GetByID(1)
		assert.NoError(t, err)
		assert.Equal(t, expectedPR, pr)
		mockRepo.AssertExpectations(t)
	})
}

func TestPullRequestService_Update(t *testing.T) {
	t.Run("successful PR update", func(t *testing.T) {
		mockRepo := new(MockPullRequestRepository)
		prService := impl.NewPullRequestService(mockRepo)

		author := &models.User{
			ID:       1,
			Name:     "John Doe",
			Email:    "john@example.com",
			TeamName: "backend",
			IsActive: true,
		}

		pr := &models.PullRequest{
			ID:        1,
			Name:      "Updated feature",
			Status:    models.StatusOpen,
			Author:    author,
			Reviewers: []*models.User{},
			CreatedAt: time.Now(),
		}

		mockRepo.On("Update", pr).Return(nil)

		err := prService.Update(pr)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestPullRequestService_ReassignReviewers(t *testing.T) {
	t.Run("no alternative reviewers available", func(t *testing.T) {
		mockRepo := new(MockPullRequestRepository)
		prService := impl.NewPullRequestService(mockRepo)

		author := &models.User{
			ID:       1,
			Name:     "John Doe",
			Email:    "john@example.com",
			TeamName: "backend",
			IsActive: true,
		}

		oldReviewer := &models.User{
			ID:       2,
			Name:     "Only Reviewer",
			Email:    "only@example.com",
			TeamName: "backend",
			IsActive: true,
		}

		pr := &models.PullRequest{
			ID:        1,
			Name:      "Feature PR",
			Status:    models.StatusOpen,
			Author:    author,
			Reviewers: []*models.User{oldReviewer},
			CreatedAt: time.Now(),
		}

		existingPR := &models.PullRequest{
			ID:        1,
			Name:      "Feature PR",
			Status:    models.StatusOpen,
			Author:    author,
			Reviewers: []*models.User{oldReviewer},
			CreatedAt: time.Now(),
		}

		possibleReviewers := []*models.User{
			{ID: 2, Name: "Only Reviewer", Email: "only@example.com", TeamName: "backend", IsActive: true},
		}

		mockRepo.On("GetByID", 1).Return(existingPR, nil)
		mockRepo.On("FindPossibleReviewers", author).Return(possibleReviewers, nil)

		err := prService.ReassignReviewers(pr, oldReviewer)
		assert.ErrorIs(t, err, models.ErrReviewerNotFound)
		mockRepo.AssertExpectations(t)
	})
}

func TestPullRequestService_MergeRequest(t *testing.T) {
	t.Run("successful merge request", func(t *testing.T) {
		mockRepo := new(MockPullRequestRepository)
		prService := impl.NewPullRequestService(mockRepo)

		author := &models.User{
			ID:       1,
			Name:     "John Doe",
			Email:    "john@example.com",
			TeamName: "backend",
			IsActive: true,
		}

		pr := &models.PullRequest{
			ID:        1,
			Name:      "Feature PR",
			Status:    models.StatusOpen,
			Author:    author,
			Reviewers: []*models.User{},
			CreatedAt: time.Now(),
		}

		existingPR := &models.PullRequest{
			ID:        1,
			Name:      "Feature PR",
			Status:    models.StatusOpen,
			Author:    author,
			Reviewers: []*models.User{},
			CreatedAt: time.Now(),
		}

		mockRepo.On("GetByID", 1).Return(existingPR, nil)
		mockRepo.On("Update", mock.MatchedBy(func(pr *models.PullRequest) bool {
			return pr.Status == models.StatusMerged && !pr.MergedAt.IsZero()
		})).Return(nil)

		err := prService.MergeRequest(pr)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("PR not found for merge", func(t *testing.T) {
		mockRepo := new(MockPullRequestRepository)
		prService := impl.NewPullRequestService(mockRepo)

		author := &models.User{
			ID:       1,
			Name:     "John Doe",
			Email:    "john@example.com",
			TeamName: "backend",
			IsActive: true,
		}

		pr := &models.PullRequest{
			ID:        999,
			Name:      "Nonexistent PR",
			Status:    models.StatusOpen,
			Author:    author,
			Reviewers: []*models.User{},
			CreatedAt: time.Now(),
		}

		mockRepo.On("GetByID", 999).Return(nil, repositories.ErrPullRequestNotFoundInPersistence)

		err := prService.MergeRequest(pr)
		assert.ErrorIs(t, err, repositories.ErrPullRequestNotFoundInPersistence)
		mockRepo.AssertExpectations(t)
	})
}
