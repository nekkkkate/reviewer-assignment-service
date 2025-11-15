package services

import (
	"reviewer-assignment-service/internal/domain/models"
)

type PullRequestService interface {
	Create(pr *models.PullRequest) error
	GetByID(id int) (*models.PullRequest, error)
	GetByAuthorID(authorID int) ([]*models.PullRequest, error)
	GetByReviewerID(reviewerID int) ([]*models.PullRequest, error)
	Update(pr *models.PullRequest) error
	ReassignReviewers(pr *models.PullRequest, oldReviewer *models.User) error
	MergeRequest(pr *models.PullRequest) error
}

type UserService interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetAll() ([]*models.User, error)
	Update(user *models.User) error
	SetActive(userID int, isActive bool) error
	Deactivate(userID int) error
}

type TeamService interface {
	Create(team *models.Team) error
	GetByID(id int) (*models.Team, error)
	GetByName(name string) (*models.Team, error)
	GetAll() ([]*models.Team, error)
	Update(team *models.Team) error
}
