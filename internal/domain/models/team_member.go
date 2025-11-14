package models

type TeamMember struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

func NewTeamMember(userID int, username string, isActive bool) *TeamMember {
	return &TeamMember{
		UserID:   userID,
		Username: username,
		IsActive: isActive,
	}
}

func (tm *TeamMember) UpdateUsername(username string) {
	tm.Username = username
}

func (tm *TeamMember) UpdateIsActive(isActive bool) {
	tm.IsActive = isActive
}
