package repositories

import (
	"errors"
	"reviewer-assignment-service/internal/domain/models"
)

type TeamRepository interface {
	Add(team *models.Team) error
	GetByID(id int) (*models.Team, error)
	GetByName(name string) (*models.Team, error)
	GetAll() ([]*models.Team, error)
	Update(team *models.Team) error
	AddUserToTeam(teamID, userID int) error
	RemoveUserFromTeam(teamID, userID int) error
}

var (
	ErrTeamNotFoundInPersistence = errors.New("team not found")
	ErrTeamAlreadyExists         = errors.New("team already exists")
)
