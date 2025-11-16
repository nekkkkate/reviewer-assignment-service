package dtos

type SetUserActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type UserResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type UserPRsResponse struct {
	UserID       string            `json:"user_id"`
	PullRequests []PRShortResponse `json:"pull_requests"`
}

type PRShortResponse struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type GetUserByEmailRequest struct {
	Email string `json:"email"`
}

type DeactivateUserRequest struct {
	UserID string `json:"user_id"`
}
