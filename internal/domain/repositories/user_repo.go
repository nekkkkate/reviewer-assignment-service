package repositories

import (
	"errors"
	"reviewer-assignment-service/internal/domain/models"
)

type UserRepository interface {
	Add(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetAll() ([]*models.User, error)
	GetActiveUsers() ([]*models.User, error)
	Update(user *models.User) error
	Deactivate(userID int) error
}

var (
	ErrUserNotFoundInPersistence = errors.New("user not found")
	ErrUserWithThatEmailNotFound = errors.New("user with that email not found")
	ErrUserAlreadyExists         = errors.New("user already exists")
)
