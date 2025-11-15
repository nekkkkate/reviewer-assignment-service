package impl

import (
	"reviewer-assignment-service/internal/domain/models"
	"reviewer-assignment-service/internal/domain/repositories"
)

type TeamService struct {
	teamRepository repositories.TeamRepository
}

func NewTeamService(teamRepository repositories.TeamRepository) *TeamService {
	return &TeamService{
		teamRepository: teamRepository,
	}
}

func (t *TeamService) Create(team *models.Team) error {
	return t.teamRepository.Add(team)
}

func (t *TeamService) GetByID(id int) (*models.Team, error) {
	return t.teamRepository.GetByID(id)
}

func (t *TeamService) GetByName(name string) (*models.Team, error) {
	return t.teamRepository.GetByName(name)
}

func (t *TeamService) GetAll() ([]*models.Team, error) {
	return t.teamRepository.GetAll()
}

func (t *TeamService) Update(team *models.Team) error {
	return t.teamRepository.Update(team)
}
