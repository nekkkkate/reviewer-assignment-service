package impl

import (
	"reviewer-assignment-service/internal/domain/models"
	"reviewer-assignment-service/internal/domain/repositories"
)

type TeamServiceImpl struct {
	teamRepository repositories.TeamRepository
}

func NewTeamService(teamRepository repositories.TeamRepository) *TeamServiceImpl {
	return &TeamServiceImpl{
		teamRepository: teamRepository,
	}
}

func (t *TeamServiceImpl) Create(team *models.Team) error {
	return t.teamRepository.Add(team)
}

func (t *TeamServiceImpl) GetByID(id int) (*models.Team, error) {
	return t.teamRepository.GetByID(id)
}

func (t *TeamServiceImpl) GetByName(name string) (*models.Team, error) {
	return t.teamRepository.GetByName(name)
}

func (t *TeamServiceImpl) GetAll() ([]*models.Team, error) {
	return t.teamRepository.GetAll()
}

func (t *TeamServiceImpl) Update(team *models.Team) error {
	return t.teamRepository.Update(team)
}
