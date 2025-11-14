package models

import (
	"errors"
	"time"
)

type PullRequest struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Status    PRStatus  `json:"status"`
	Author    *User     `json:"author"`
	Reviewers []*User   `json:"reviewers"`
	CreatedAt time.Time `json:"created_at"`
	MergedAt  time.Time `json:"merged_at"`
}

type PRStatus string

const (
	StatusOpen   PRStatus = "OPEN"
	StatusMerged PRStatus = "MERGED"
)

func NewPullRequest(name string, author *User, team *Team) (*PullRequest, error) {
	if !team.IsMemberInTeam(author.ID) {
		return nil, ErrAuthorNotInTeam
	}
	return &PullRequest{
		Name:      name,
		Status:    StatusOpen,
		Author:    author,
		Reviewers: make([]*User, 0, 2),
		CreatedAt: time.Now(),
	}, nil
}

func (pr *PullRequest) CanModifyReviewers() bool {
	return pr.Status == StatusOpen
}

func (pr *PullRequest) SetStatusMerged() {
	pr.Status = StatusMerged
}
func (pr *PullRequest) SetMergedAt(mergedAt time.Time) {
	pr.MergedAt = mergedAt
}

func (pr *PullRequest) SetId(id int) {
	pr.ID = id
}

func (pr *PullRequest) AddReviewer(reviewer *User) error {
	if !pr.CanModifyReviewers() {
		return ErrPRAlreadyMerged
	}

	if len(pr.Reviewers) >= 2 {
		return ErrTooManyReviewers
	}

	for _, r := range pr.Reviewers {
		if r.ID == reviewer.ID {
			return ErrReviewerAlreadyAssigned
		}
	}

	pr.Reviewers = append(pr.Reviewers, reviewer)
	return nil
}

func (pr *PullRequest) RemoveReviewer(reviewerID int) error {
	if !pr.CanModifyReviewers() {
		return ErrPRAlreadyMerged
	}

	for i, reviewer := range pr.Reviewers {
		if reviewer.ID == reviewerID {
			pr.Reviewers = append(pr.Reviewers[:i], pr.Reviewers[i+1:]...)
			return nil
		}
	}
	return ErrReviewerNotFound
}

func (pr *PullRequest) ReplaceReviewer(oldReviewerID int, newReviewer *User) error {
	if err := pr.RemoveReviewer(oldReviewerID); err != nil {
		return err
	}
	return pr.AddReviewer(newReviewer)
}

var (
	ErrAuthorNotInTeam         = errors.New("author not in team")
	ErrReviewerNotFound        = errors.New("reviewer not found")
	ErrPRAlreadyMerged         = errors.New("pull request already merged")
	ErrReviewerAlreadyAssigned = errors.New("reviewer already assigned")
	ErrTooManyReviewers        = errors.New("too many reviewers")
)
