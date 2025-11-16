package impl

import (
	"reviewer-assignment-service/internal/domain/models"
	"reviewer-assignment-service/internal/domain/repositories"
	"time"
)

type PullRequestServiceImpl struct {
	pullRequestRepository repositories.PullRequestRepository
}

func NewPullRequestService(pullRequestRepository repositories.PullRequestRepository) *PullRequestServiceImpl {
	return &PullRequestServiceImpl{
		pullRequestRepository: pullRequestRepository,
	}
}

func (p *PullRequestServiceImpl) Create(pr *models.PullRequest) error {
	return p.pullRequestRepository.Add(pr)
}

func (p *PullRequestServiceImpl) GetByID(id int) (*models.PullRequest, error) {
	return p.pullRequestRepository.GetByID(id)
}

func (p *PullRequestServiceImpl) Update(pr *models.PullRequest) error {
	return p.pullRequestRepository.Update(pr)
}

func (p *PullRequestServiceImpl) ReassignReviewers(pr *models.PullRequest, oldReviewer *models.User) error {
	pullRequest, err := p.pullRequestRepository.GetByID(pr.ID)
	if err != nil {
		return err
	}
	possibleReviewers, err := p.pullRequestRepository.FindPossibleReviewers(pr.Author)
	if err != nil {
		return err
	}
	newReviewerID := defaultId
	newReviewer := &models.User{}
	for _, reviewer := range possibleReviewers {
		if reviewer.ID != oldReviewer.ID {
			newReviewerID = reviewer.ID
			newReviewer = reviewer
			break
		}
	}
	if newReviewerID == defaultId {
		return models.ErrReviewerNotFound
	}
	return pullRequest.ReplaceReviewer(oldReviewer.ID, newReviewer)
}

func (p *PullRequestServiceImpl) MergeRequest(pr *models.PullRequest) error {
	pullRequest, err := p.pullRequestRepository.GetByID(pr.ID)
	if err != nil {
		return err
	}
	pullRequest.Status = models.StatusMerged
	pullRequest.SetMergedAt(time.Now())
	return p.pullRequestRepository.Update(pullRequest)
}

var defaultId = -10

func (p *PullRequestServiceImpl) GetByAuthorID(authorID int) ([]*models.PullRequest, error) {
	return p.pullRequestRepository.GetByAuthorID(authorID)
}
func (p *PullRequestServiceImpl) GetByReviewerID(reviewerID int) ([]*models.PullRequest, error) {
	return p.pullRequestRepository.GetByReviewerID(reviewerID)
}
