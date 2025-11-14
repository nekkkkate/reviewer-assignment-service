package in_memory

import (
	"reviewer-assignment-service/internal/domain/models"
	"reviewer-assignment-service/internal/domain/repositories"
)

type UserRepository struct {
	users map[int]*models.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[int]*models.User),
	}
}

func (r *UserRepository) Add(user *models.User) error {
	if _, ok := r.users[user.ID]; ok {
		return repositories.ErrUserAlreadyExists
	}
	r.users[user.ID] = user
	return nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	if user, ok := r.users[id]; ok {
		return user, nil
	}
	return nil, repositories.ErrUserNotFoundInPersistence
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, repositories.ErrUserWithThatEmailNotFound
}

func (r *UserRepository) GetAll() ([]*models.User, error) {
	var users []*models.User
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) GetActiveUsers() ([]*models.User, error) {
	var activeUsers []*models.User
	for _, user := range r.users {
		if user.IsActive {
			activeUsers = append(activeUsers, user)
		}
	}
	return activeUsers, nil
}

func (r *UserRepository) Update(user *models.User) error {
	if _, ok := r.users[user.ID]; !ok {
		return repositories.ErrUserNotFoundInPersistence
	}
	r.users[user.ID] = user
	return nil
}
