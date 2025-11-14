package in_memory

import (
	"reviewer-assignment-service/internal/domain/models"
	"reviewer-assignment-service/internal/domain/repositories"
)

type TeamRepository struct {
	teams map[int]*models.Team
}

func NewTeamRepository() *TeamRepository {
	return &TeamRepository{
		teams: make(map[int]*models.Team),
	}
}

func (r *TeamRepository) Add(team *models.Team) error {
	if _, ok := r.teams[team.ID]; ok {
		return repositories.ErrTeamAlreadyExists
	}
	r.teams[team.ID] = team
	return nil
}
func (r *TeamRepository) GetByID(id int) (*models.Team, error) {
	if team, ok := r.teams[id]; ok {
		return team, nil
	}
	return nil, repositories.ErrTeamNotFoundInPersistence
}

func (r *TeamRepository) GetByName(name string) (*models.Team, error) {
	for _, team := range r.teams {
		if team.Name == name {
			return team, nil
		}
	}
	return nil, repositories.ErrTeamNotFoundInPersistence
}

func (r *TeamRepository) GetAll() ([]*models.Team, error) {
	var teams []*models.Team
	for _, team := range r.teams {
		teams = append(teams, team)
	}
	return teams, nil
}

func (r *TeamRepository) Update(team *models.Team) error {
	if _, ok := r.teams[team.ID]; !ok {
		return repositories.ErrTeamNotFoundInPersistence
	}
	r.teams[team.ID] = team
	return nil
}
