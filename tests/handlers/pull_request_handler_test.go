package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	appHandlers "reviewer-assignment-service/internal/app/handlers"
	"reviewer-assignment-service/internal/app/transport/dtos"
	"reviewer-assignment-service/internal/domain/models"
	"reviewer-assignment-service/internal/domain/services"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockPullRequestService struct {
	mock.Mock
}

func (m *MockPullRequestService) Create(pr *models.PullRequest) error {
	args := m.Called(pr)
	return args.Error(0)
}

func (m *MockPullRequestService) GetByID(id int) (*models.PullRequest, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PullRequest), args.Error(1)
}

func (m *MockPullRequestService) GetByAuthorID(authorID int) ([]*models.PullRequest, error) {
	args := m.Called(authorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.PullRequest), args.Error(1)
}

func (m *MockPullRequestService) GetByReviewerID(reviewerID int) ([]*models.PullRequest, error) {
	args := m.Called(reviewerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.PullRequest), args.Error(1)
}

func (m *MockPullRequestService) Update(pr *models.PullRequest) error {
	args := m.Called(pr)
	return args.Error(0)
}

func (m *MockPullRequestService) ReassignReviewers(pr *models.PullRequest, oldReviewer *models.User) error {
	args := m.Called(pr, oldReviewer)
	return args.Error(0)
}

func (m *MockPullRequestService) MergeRequest(pr *models.PullRequest) error {
	args := m.Called(pr)
	return args.Error(0)
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) GetByID(id int) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetAll() ([]*models.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserService) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) SetActive(userID int, isActive bool) error {
	args := m.Called(userID, isActive)
	return args.Error(0)
}

func (m *MockUserService) Deactivate(userID int) error {
	args := m.Called(userID)
	return args.Error(0)
}

var _ services.PullRequestService = (*MockPullRequestService)(nil)
var _ services.UserService = (*MockUserService)(nil)

func TestPullRequestHandler_CreatePullRequest_Success(t *testing.T) {
	mockPRService := new(MockPullRequestService)
	mockUserService := new(MockUserService)
	handler := appHandlers.NewPullRequestHandler(mockPRService, mockUserService)

	author := &models.User{ID: 1, Name: "Author", Email: "author@example.com", TeamName: "backend", IsActive: true}
	reviewer1 := &models.User{ID: 2, Name: "Reviewer1", Email: "rev1@example.com", TeamName: "backend", IsActive: true}
	reviewer2 := &models.User{ID: 3, Name: "Reviewer2", Email: "rev2@example.com", TeamName: "backend", IsActive: true}

	mockUserService.On("GetByID", 1).Return(author, nil)
	mockUserService.On("GetByID", 2).Return(reviewer1, nil)
	mockUserService.On("GetByID", 3).Return(reviewer2, nil)

	mockPRService.On("Create", mock.MatchedBy(func(pr *models.PullRequest) bool {
		return pr.Name == "New PR" &&
			pr.Status == models.StatusOpen &&
			pr.Author == author &&
			len(pr.Reviewers) == 2 &&
			pr.Reviewers[0] == reviewer1 &&
			pr.Reviewers[1] == reviewer2
	})).Return(nil)

	reqBody := dtos.CreatePullRequestRequest{
		Name:      "New PR",
		AuthorID:  1,
		Reviewers: []int{2, 3},
	}
	bodyBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/pull-requests", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()

	handler.CreatePullRequest(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp dtos.PullRequestResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, "New PR", resp.Name)
	assert.Equal(t, string(models.StatusOpen), resp.Status)
	if assert.NotNil(t, resp.Author) {
		assert.Equal(t, "1", resp.Author.UserID)
	}
	if assert.Len(t, resp.Reviewers, 2) {
		assert.Equal(t, "2", resp.Reviewers[0].UserID)
		assert.Equal(t, "3", resp.Reviewers[1].UserID)
	}
	assert.False(t, resp.CreatedAt.IsZero())
	assert.Nil(t, resp.MergedAt)

	mockUserService.AssertExpectations(t)
	mockPRService.AssertExpectations(t)
}

func TestPullRequestHandler_GetPullRequestByID_Success(t *testing.T) {
	mockPRService := new(MockPullRequestService)
	mockUserService := new(MockUserService)
	handler := appHandlers.NewPullRequestHandler(mockPRService, mockUserService)

	author := &models.User{ID: 1, Name: "Author", Email: "author@example.com", TeamName: "backend", IsActive: true}
	reviewer := &models.User{ID: 2, Name: "Reviewer", Email: "rev@example.com", TeamName: "backend", IsActive: true}
	createdAt := time.Now()

	pr := &models.PullRequest{
		ID:        1,
		Name:      "Existing PR",
		Status:    models.StatusOpen,
		Author:    author,
		Reviewers: []*models.User{reviewer},
		CreatedAt: createdAt,
	}

	mockPRService.On("GetByID", 1).Return(pr, nil)

	req := httptest.NewRequest(http.MethodGet, "/pull-requests/1", nil)
	rec := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetPullRequestByID(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dtos.PullRequestResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, 1, resp.ID)
	assert.Equal(t, "Existing PR", resp.Name)
	assert.Equal(t, string(models.StatusOpen), resp.Status)
	if assert.NotNil(t, resp.Author) {
		assert.Equal(t, "1", resp.Author.UserID)
	}
	if assert.Len(t, resp.Reviewers, 1) {
		assert.Equal(t, "2", resp.Reviewers[0].UserID)
	}

	mockPRService.AssertExpectations(t)
}

func TestPullRequestHandler_UpdatePullRequest_Success(t *testing.T) {
	mockPRService := new(MockPullRequestService)
	mockUserService := new(MockUserService)
	handler := appHandlers.NewPullRequestHandler(mockPRService, mockUserService)

	author := &models.User{ID: 1, Name: "Author", Email: "author@example.com", TeamName: "backend", IsActive: true}
	reviewer := &models.User{ID: 2, Name: "Reviewer", Email: "rev@example.com", TeamName: "backend", IsActive: true}

	existingPR := &models.PullRequest{
		ID:        1,
		Name:      "Old PR",
		Status:    models.StatusOpen,
		Author:    author,
		Reviewers: []*models.User{},
		CreatedAt: time.Now(),
	}

	mockPRService.On("GetByID", 1).Return(existingPR, nil)
	mockUserService.On("GetByID", 2).Return(reviewer, nil)
	mockPRService.On("Update", mock.MatchedBy(func(pr *models.PullRequest) bool {
		return pr.ID == 1 &&
			pr.Name == "Updated PR" &&
			pr.Status == models.StatusOpen &&
			len(pr.Reviewers) == 1 &&
			pr.Reviewers[0] == reviewer
	})).Return(nil)

	reqBody := dtos.UpdatePullRequestRequest{
		Name:      "Updated PR",
		Status:    "OPEN",
		Reviewers: []int{2},
	}
	bodyBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/pull-requests/1", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.UpdatePullRequest(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dtos.PullRequestResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, 1, resp.ID)
	assert.Equal(t, "Updated PR", resp.Name)
	assert.Equal(t, "OPEN", resp.Status)
	if assert.Len(t, resp.Reviewers, 1) {
		assert.Equal(t, "2", resp.Reviewers[0].UserID)
	}

	mockPRService.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
}
