package postgres

import (
	"database/sql"
	"reviewer-assignment-service/internal/domain/models"

	"github.com/Masterminds/squirrel"
)

type TeamDataBase struct {
	db *sql.DB
	sb squirrel.StatementBuilderType
}

func (t TeamDataBase) Add(team *models.Team) error {
	//TODO implement me
	panic("implement me")
}

func (t TeamDataBase) GetByID(id int) (*models.Team, error) {
	//TODO implement me
	panic("implement me")
}

func (t TeamDataBase) GetByName(name string) (*models.Team, error) {
	//TODO implement me
	panic("implement me")
}

func (t TeamDataBase) GetAll() ([]*models.Team, error) {
	//TODO implement me
	panic("implement me")
}

func (t TeamDataBase) Update(team *models.Team) error {
	//TODO implement me
	panic("implement me")
}

func (t TeamDataBase) AddUserToTeam(teamID, userID int) error {
	//TODO implement me
	panic("implement me")
}

func (t TeamDataBase) RemoveUserFromTeam(teamID, userID int) error {
	//TODO implement me
	panic("implement me")
}

func NewTeamDataBase(db *sql.DB) *TeamDataBase {
	return &TeamDataBase{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
