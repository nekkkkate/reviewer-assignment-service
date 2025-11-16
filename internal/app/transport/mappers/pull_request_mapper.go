package mappers

import (
	"reviewer-assignment-service/internal/app/transport/dtos"
	"reviewer-assignment-service/internal/domain/models"
	"time"
)

func ToPullRequestResponse(pr *models.PullRequest) *dtos.PullRequestResponse {
	authorResponse := UserToResponse(pr.Author)
	response := &dtos.PullRequestResponse{
		ID:        pr.ID,
		Name:      pr.Name,
		Status:    string(pr.Status),
		Author:    &authorResponse,
		Reviewers: make([]*dtos.UserResponse, len(pr.Reviewers)),
		CreatedAt: pr.CreatedAt,
	}

	if !pr.MergedAt.IsZero() {
		response.MergedAt = &pr.MergedAt
	}

	for i, reviewer := range pr.Reviewers {
		reviewerResponse := UserToResponse(reviewer)
		response.Reviewers[i] = &reviewerResponse
	}

	return response
}

func ToPullRequestListResponse(prs []*models.PullRequest) *dtos.PullRequestListResponse {
	responses := make([]*dtos.PullRequestResponse, len(prs))
	for i, pr := range prs {
		responses[i] = ToPullRequestResponse(pr)
	}

	return &dtos.PullRequestListResponse{
		PullRequests: responses,
		Total:        len(prs),
	}
}

func ToPullRequestModel(req *dtos.CreatePullRequestRequest, author *models.User) *models.PullRequest {
	pr := &models.PullRequest{
		Name:      req.Name,
		Status:    models.StatusOpen,
		Author:    author,
		Reviewers: make([]*models.User, 0),
		CreatedAt: time.Now(),
	}

	return pr
}

func UpdatePullRequestFromRequest(pr *models.PullRequest, req *dtos.UpdatePullRequestRequest) {
	pr.Name = req.Name
	pr.Status = models.PRStatus(req.Status)
}
