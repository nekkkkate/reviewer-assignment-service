package service

import (
	"reviewer-assignment-service/internal/domain/models"
	"reviewer-assignment-service/internal/domain/repositories"
	"reviewer-assignment-service/internal/domain/services/impl"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTeamRepository struct {
	mock.Mock
}

func (m *MockTeamRepository) Add(team *models.Team) error {
	args := m.Called(team)
	return args.Error(0)
}

func (m *MockTeamRepository) GetByID(id int) (*models.Team, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Team), args.Error(1)
}

func (m *MockTeamRepository) GetByName(name string) (*models.Team, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Team), args.Error(1)
}

func (m *MockTeamRepository) GetAll() ([]*models.Team, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Team), args.Error(1)
}

func (m *MockTeamRepository) Update(team *models.Team) error {
	args := m.Called(team)
	return args.Error(0)
}

func (m *MockTeamRepository) AddUserToTeam(teamID, userID int) error {
	args := m.Called(teamID, userID)
	return args.Error(0)
}

func (m *MockTeamRepository) RemoveUserFromTeam(teamID, userID int) error {
	args := m.Called(teamID, userID)
	return args.Error(0)
}

func TestTeamService_Create(t *testing.T) {
	t.Run("successful team creation", func(t *testing.T) {
		mockRepo := new(MockTeamRepository)
		teamService := impl.NewTeamService(mockRepo)

		team := &models.Team{
			Name:    "backend",
			Members: make(map[int]*models.TeamMember),
		}

		mockRepo.On("Add", team).Return(nil)

		err := teamService.Create(team)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("team creation with duplicate name", func(t *testing.T) {
		mockRepo := new(MockTeamRepository)
		teamService := impl.NewTeamService(mockRepo)

		team := &models.Team{
			Name:    "backend",
			Members: make(map[int]*models.TeamMember),
		}

		mockRepo.On("Add", team).Return(repositories.ErrTeamAlreadyExists)

		err := teamService.Create(team)
		assert.ErrorIs(t, err, repositories.ErrTeamAlreadyExists)
		mockRepo.AssertExpectations(t)
	})
}

func TestTeamService_GetByID(t *testing.T) {
	t.Run("successful get by id", func(t *testing.T) {
		mockRepo := new(MockTeamRepository)
		teamService := impl.NewTeamService(mockRepo)

		expectedTeam := &models.Team{
			ID:      1,
			Name:    "backend",
			Members: make(map[int]*models.TeamMember),
		}

		mockRepo.On("GetByID", 1).Return(expectedTeam, nil)

		team, err := teamService.GetByID(1)
		assert.NoError(t, err)
		assert.Equal(t, expectedTeam, team)
		mockRepo.AssertExpectations(t)
	})

	t.Run("team not found", func(t *testing.T) {
		mockRepo := new(MockTeamRepository)
		teamService := impl.NewTeamService(mockRepo)

		mockRepo.On("GetByID", 999).Return(nil, repositories.ErrTeamNotFoundInPersistence)

		team, err := teamService.GetByID(999)
		assert.Nil(t, team)
		assert.ErrorIs(t, err, repositories.ErrTeamNotFoundInPersistence)
		mockRepo.AssertExpectations(t)
	})
}

func TestTeamService_GetByName(t *testing.T) {
	t.Run("successful get by name", func(t *testing.T) {
		mockRepo := new(MockTeamRepository)
		teamService := impl.NewTeamService(mockRepo)

		expectedTeam := &models.Team{
			ID:      1,
			Name:    "backend",
			Members: make(map[int]*models.TeamMember),
		}

		mockRepo.On("GetByName", "backend").Return(expectedTeam, nil)

		team, err := teamService.GetByName("backend")
		assert.NoError(t, err)
		assert.Equal(t, expectedTeam, team)
		mockRepo.AssertExpectations(t)
	})
}

func TestTeamService_GetAll(t *testing.T) {
	t.Run("successful get all teams", func(t *testing.T) {
		mockRepo := new(MockTeamRepository)
		teamService := impl.NewTeamService(mockRepo)

		expectedTeams := []*models.Team{
			{ID: 1, Name: "backend", Members: make(map[int]*models.TeamMember)},
			{ID: 2, Name: "frontend", Members: make(map[int]*models.TeamMember)},
		}

		mockRepo.On("GetAll").Return(expectedTeams, nil)

		teams, err := teamService.GetAll()
		assert.NoError(t, err)
		assert.Equal(t, expectedTeams, teams)
		mockRepo.AssertExpectations(t)
	})
}

func TestTeamService_Update(t *testing.T) {
	t.Run("successful team update", func(t *testing.T) {
		mockRepo := new(MockTeamRepository)
		teamService := impl.NewTeamService(mockRepo)

		team := &models.Team{
			ID:      1,
			Name:    "backend-updated",
			Members: make(map[int]*models.TeamMember),
		}

		mockRepo.On("Update", team).Return(nil)

		err := teamService.Update(team)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
