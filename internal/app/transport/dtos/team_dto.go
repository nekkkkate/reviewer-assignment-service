package dtos

type CreateTeamRequest struct {
	Name    string                    `json:"name" binding:"required,min=2,max=100"`
	Members []CreateTeamMemberRequest `json:"members,omitempty"`
}

type CreateTeamMemberRequest struct {
	UserID   int    `json:"user_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	IsActive bool   `json:"is_active"`
}

type UpdateTeamRequest struct {
	Name    string                    `json:"name" binding:"required,min=2,max=100"`
	Members []CreateTeamMemberRequest `json:"members,omitempty"`
}

type AddMemberRequest struct {
	UserID int `json:"user_id" binding:"required"`
}

type TeamResponse struct {
	ID      int                  `json:"id"`
	Name    string               `json:"name"`
	Members []TeamMemberResponse `json:"members"`
}

type TeamMemberResponse struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type TeamListResponse struct {
	Teams []TeamResponse `json:"teams"`
	Total int            `json:"total"`
}
