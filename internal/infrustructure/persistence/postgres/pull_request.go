package postgres

import (
	"database/sql"
	"reviewer-assignment-service/internal/domain/models"

	"github.com/Masterminds/squirrel"
)

type PullRequestDataBase struct {
	db *sql.DB
	sb squirrel.StatementBuilderType
}

func (p PullRequestDataBase) Add(pr *models.PullRequest) error {
	//TODO implement me
	panic("implement me")
}

func (p PullRequestDataBase) GetByID(id int) (*models.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (p PullRequestDataBase) GetAll() ([]*models.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (p PullRequestDataBase) GetByStatus(status models.PRStatus) ([]*models.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (p PullRequestDataBase) GetByAuthorID(authorID int) ([]*models.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (p PullRequestDataBase) GetByReviewerID(reviewerID int) ([]*models.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (p PullRequestDataBase) Update(pr *models.PullRequest) error {
	//TODO implement me
	panic("implement me")
}

func NewPullRequestDataBase(db *sql.DB) *PullRequestDataBase {
	return &PullRequestDataBase{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
