package service

import (
	"reviewer-assignment-service/internal/domain/models"
	"reviewer-assignment-service/internal/domain/repositories"
	"reviewer-assignment-service/internal/domain/services/impl"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Add(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id int) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetAll() ([]*models.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Deactivate(userID int) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserRepository) GetWithFilters(teamName string, isActive bool) ([]*models.User, error) {
	args := m.Called(teamName, isActive)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRepository) GetActiveUsers() ([]*models.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func TestUserService_Create(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userService := impl.NewUserService(mockRepo)

		user := &models.User{
			Name:     "John Doe",
			Email:    "john@example.com",
			TeamName: "backend",
			IsActive: true,
		}

		mockRepo.On("Add", user).Return(nil)

		err := userService.Create(user)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("creation with duplicate email", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userService := impl.NewUserService(mockRepo)

		user := &models.User{
			Name:     "John Doe",
			Email:    "john@example.com",
			TeamName: "backend",
			IsActive: true,
		}

		mockRepo.On("Add", user).Return(repositories.ErrUserAlreadyExists)

		err := userService.Create(user)
		assert.ErrorIs(t, err, repositories.ErrUserAlreadyExists)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetByID(t *testing.T) {
	t.Run("successful get by id", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userService := impl.NewUserService(mockRepo)

		expectedUser := &models.User{
			ID:       1,
			Name:     "John Doe",
			Email:    "john@example.com",
			TeamName: "backend",
			IsActive: true,
		}

		mockRepo.On("GetByID", 1).Return(expectedUser, nil)

		user, err := userService.GetByID(1)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userService := impl.NewUserService(mockRepo)

		mockRepo.On("GetByID", 999).Return(nil, repositories.ErrUserNotFoundInPersistence)

		user, err := userService.GetByID(999)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, repositories.ErrUserNotFoundInPersistence)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetByEmail(t *testing.T) {
	t.Run("successful get by email", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userService := impl.NewUserService(mockRepo)

		expectedUser := &models.User{
			ID:       1,
			Name:     "John Doe",
			Email:    "john@example.com",
			TeamName: "backend",
			IsActive: true,
		}

		mockRepo.On("GetByEmail", "john@example.com").Return(expectedUser, nil)

		user, err := userService.GetByEmail("john@example.com")
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetAll(t *testing.T) {
	t.Run("successful get all users", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userService := impl.NewUserService(mockRepo)

		expectedUsers := []*models.User{
			{ID: 1, Name: "John Doe", Email: "john@example.com", TeamName: "backend", IsActive: true},
			{ID: 2, Name: "Jane Smith", Email: "jane@example.com", TeamName: "frontend", IsActive: true},
		}

		mockRepo.On("GetAll").Return(expectedUsers, nil)

		users, err := userService.GetAll()
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Update(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userService := impl.NewUserService(mockRepo)

		user := &models.User{
			ID:       1,
			Name:     "John Doe Updated",
			Email:    "john.updated@example.com",
			TeamName: "backend",
			IsActive: true,
		}

		mockRepo.On("Update", user).Return(nil)

		err := userService.Update(user)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_SetActive(t *testing.T) {
	t.Run("successful activate user", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userService := impl.NewUserService(mockRepo)

		user := &models.User{
			ID:       1,
			Name:     "John Doe",
			Email:    "john@example.com",
			TeamName: "backend",
			IsActive: false,
		}

		mockRepo.On("GetByID", 1).Return(user, nil)
		mockRepo.On("Update", mock.MatchedBy(func(u *models.User) bool {
			return u.ID == 1 && u.IsActive == true
		})).Return(nil)

		err := userService.SetActive(1, true)
		assert.NoError(t, err)
		assert.True(t, user.IsActive)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found for activation", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userService := impl.NewUserService(mockRepo)

		mockRepo.On("GetByID", 999).Return(nil, repositories.ErrUserNotFoundInPersistence)

		err := userService.SetActive(999, true)
		assert.ErrorIs(t, err, repositories.ErrUserNotFoundInPersistence)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Deactivate(t *testing.T) {
	t.Run("successful deactivate", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userService := impl.NewUserService(mockRepo)

		mockRepo.On("Deactivate", 1).Return(nil)

		err := userService.Deactivate(1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
