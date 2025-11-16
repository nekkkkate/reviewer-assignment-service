package mappers

import (
	"reviewer-assignment-service/internal/app/transport/dtos"
	"reviewer-assignment-service/internal/domain/models"
)

func ToTeamResponse(team *models.Team) dtos.TeamResponse {
	memberResponses := make([]dtos.TeamMemberResponse, 0, len(team.Members))
	for _, member := range team.Members {
		memberResponses = append(memberResponses, dtos.TeamMemberResponse{
			UserID:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		})
	}

	return dtos.TeamResponse{
		ID:      team.ID,
		Name:    team.Name,
		Members: memberResponses,
	}
}

func ToTeamListResponse(teams []*models.Team) dtos.TeamListResponse {
	teamResponses := make([]dtos.TeamResponse, len(teams))
	for i, team := range teams {
		teamResponses[i] = ToTeamResponse(team)
	}

	return dtos.TeamListResponse{
		Teams: teamResponses,
		Total: len(teams),
	}
}

func ToTeamModel(req dtos.CreateTeamRequest) *models.Team {
	team := models.NewTeam(req.Name)

	for _, memberReq := range req.Members {
		member := models.NewTeamMember(memberReq.UserID, memberReq.Username, memberReq.IsActive)
		team.Members[member.UserID] = member
	}

	return team
}

func UpdateTeamFromRequest(team *models.Team, req dtos.UpdateTeamRequest) {
	team.UpdateName(req.Name)

	team.Members = make(map[int]*models.TeamMember)
	for _, memberReq := range req.Members {
		member := models.NewTeamMember(memberReq.UserID, memberReq.Username, memberReq.IsActive)
		team.Members[member.UserID] = member
	}
}
