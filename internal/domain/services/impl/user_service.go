package impl

import (
	"reviewer-assignment-service/internal/domain/models"
	"reviewer-assignment-service/internal/domain/repositories"
)

type UserService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (u *UserService) Create(user *models.User) error {
	return u.userRepository.Add(user)
}

func (u *UserService) GetByID(id int) (*models.User, error) {
	return u.userRepository.GetByID(id)
}

func (u *UserService) GetByEmail(email string) (*models.User, error) {
	return u.userRepository.GetByEmail(email)
}

func (u *UserService) GetAll() ([]*models.User, error) {
	return u.userRepository.GetAll()
}

func (u *UserService) Update(user *models.User) error {
	return u.userRepository.Update(user)
}

func (u *UserService) SetActive(userID int, isActive bool) error {
	user, err := u.userRepository.GetByID(userID)
	if err != nil {
		return err
	}
	user.UpdateIsActive(isActive)
	return u.userRepository.Update(user)
}
func (u *UserService) Deactivate(userID int) error {
	return u.userRepository.Deactivate(userID)
}
