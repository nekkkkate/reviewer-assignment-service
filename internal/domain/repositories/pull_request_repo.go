package repositories

import (
	"errors"
	"reviewer-assignment-service/internal/domain/models"
)

type PullRequestRepository interface {
	Add(pr *models.PullRequest) error
	GetByID(id int) (*models.PullRequest, error)
	GetAll() ([]*models.PullRequest, error)
	GetByStatus(status models.PRStatus) ([]*models.PullRequest, error)
	GetByAuthorID(authorID int) ([]*models.PullRequest, error)
	GetByReviewerID(reviewerID int) ([]*models.PullRequest, error)
	Update(pr *models.PullRequest) error
	FindPossibleReviewers(author *models.User) ([]*models.User, error)
}

var (
	ErrPullRequestNotFoundInPersistence = errors.New("pull request not found")
	ErrPullRequestAlreadyExists         = errors.New("pull request already exists")
)
