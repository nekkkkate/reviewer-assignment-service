package impl

import (
	"reviewer-assignment-service/internal/domain/models"
	"reviewer-assignment-service/internal/domain/repositories"
)

type UserServiceImpl struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepository: userRepository,
	}
}

func (u *UserServiceImpl) Create(user *models.User) error {
	return u.userRepository.Add(user)
}

func (u *UserServiceImpl) GetByID(id int) (*models.User, error) {
	return u.userRepository.GetByID(id)
}

func (u *UserServiceImpl) GetByEmail(email string) (*models.User, error) {
	return u.userRepository.GetByEmail(email)
}

func (u *UserServiceImpl) GetAll() ([]*models.User, error) {
	return u.userRepository.GetAll()
}

func (u *UserServiceImpl) Update(user *models.User) error {
	return u.userRepository.Update(user)
}

func (u *UserServiceImpl) SetActive(userID int, isActive bool) error {
	user, err := u.userRepository.GetByID(userID)
	if err != nil {
		return err
	}
	user.UpdateIsActive(isActive)
	return u.userRepository.Update(user)
}
func (u *UserServiceImpl) Deactivate(userID int) error {
	return u.userRepository.Deactivate(userID)
}
