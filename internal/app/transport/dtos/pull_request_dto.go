package dtos

import "time"

type CreatePullRequestRequest struct {
	Name      string `json:"name"`
	AuthorID  int    `json:"author_id"`
	Reviewers []int  `json:"reviewers,omitempty"`
}

type UpdatePullRequestRequest struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	Reviewers []int  `json:"reviewers,omitempty"`
}

type ReassignReviewersRequest struct {
	OldReviewerID int `json:"old_reviewer_id"`
}

type AddReviewerRequest struct {
	ReviewerID int `json:"reviewer_id"`
}

type PullRequestResponse struct {
	ID        int             `json:"id"`
	Name      string          `json:"name"`
	Status    string          `json:"status"`
	Author    *UserResponse   `json:"author"`
	Reviewers []*UserResponse `json:"reviewers"`
	CreatedAt time.Time       `json:"created_at"`
	MergedAt  *time.Time      `json:"merged_at,omitempty"`
}

type PullRequestListResponse struct {
	PullRequests []*PullRequestResponse `json:"pull_requests"`
	Total        int                    `json:"total"`
}
