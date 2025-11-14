package in_memory

import (
	"reviewer-assignment-service/internal/domain/models"
	"reviewer-assignment-service/internal/domain/repositories"
)

type PullRequestRepository struct {
	pullRequests map[int]*models.PullRequest
}

func NewPullRequestRepository() *PullRequestRepository {
	return &PullRequestRepository{
		pullRequests: make(map[int]*models.PullRequest),
	}
}
func (r *PullRequestRepository) Add(pr *models.PullRequest) error {
	if _, ok := r.pullRequests[pr.ID]; ok {
		return repositories.ErrPullRequestAlreadyExists
	}
	r.pullRequests[pr.ID] = pr
	return nil
}

func (r *PullRequestRepository) GetByID(id int) (*models.PullRequest, error) {
	if pr, ok := r.pullRequests[id]; ok {
		return pr, nil
	}
	return nil, repositories.ErrPullRequestNotFoundInPersistence
}

func (r *PullRequestRepository) GetAll() ([]*models.PullRequest, error) {
	var pullRequests []*models.PullRequest
	for _, pr := range r.pullRequests {
		pullRequests = append(pullRequests, pr)
	}
	return pullRequests, nil
}

func (r *PullRequestRepository) GetByStatus(status models.PRStatus) ([]*models.PullRequest, error) {
	var pullRequests []*models.PullRequest
	for _, pr := range r.pullRequests {
		if pr.Status == status {
			pullRequests = append(pullRequests, pr)
		}
	}
	return pullRequests, nil
}

func (r *PullRequestRepository) GetByAuthorID(authorID int) ([]*models.PullRequest, error) {
	var pullRequests []*models.PullRequest
	for _, pr := range r.pullRequests {
		if pr.Author.ID == authorID {
			pullRequests = append(pullRequests, pr)
		}
	}
	return pullRequests, nil
}

func (r *PullRequestRepository) GetByReviewerID(reviewerID int) ([]*models.PullRequest, error) {
	var pullRequests []*models.PullRequest
	for _, pr := range r.pullRequests {
		for _, reviewer := range pr.Reviewers {
			if reviewer.ID == reviewerID {
				pullRequests = append(pullRequests, pr)
				break
			}
		}
	}
	return pullRequests, nil
}
