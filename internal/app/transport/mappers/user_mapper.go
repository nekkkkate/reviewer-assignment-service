package mappers

import (
	"reviewer-assignment-service/internal/app/transport/dtos"
	"reviewer-assignment-service/internal/domain/models"
	"strconv"
)

func UserToResponse(user *models.User) dtos.UserResponse {
	return dtos.UserResponse{
		UserID:   strconv.Itoa(user.ID),
		Username: user.Name,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}
}

func UserToResponseWithPRs(userID string, prs []*models.PullRequest) dtos.UserPRsResponse {
	var prResponses []dtos.PRShortResponse
	for _, pr := range prs {
		prResponses = append(prResponses, dtos.PRShortResponse{
			PullRequestID:   strconv.Itoa(pr.ID),
			PullRequestName: pr.Name,
			AuthorID:        strconv.Itoa(pr.Author.ID),
			Status:          string(pr.Status),
		})
	}

	return dtos.UserPRsResponse{
		UserID:       userID,
		PullRequests: prResponses,
	}
}

func CreateUserRequestToDomain(req dtos.CreateUserRequest) *models.User {
	user := models.NewUser(req.Username, req.Email, req.IsActive, req.TeamName)
	return user
}

func UserToDetailedResponse(user *models.User) dtos.UserResponse {
	return dtos.UserResponse{
		UserID:   strconv.Itoa(user.ID),
		Username: user.Name,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}
}
