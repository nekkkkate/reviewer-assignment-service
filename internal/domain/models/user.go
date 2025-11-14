package models

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	IsActive bool   `json:"is_active"`
	TeamName string `json:"team_name"`
}

func NewUser(name, email string, isActive bool, teamName string) *User {
	return &User{
		Name:     name,
		Email:    email,
		IsActive: isActive,
		TeamName: teamName,
	}
}

func (u *User) SetId(id int) {
	u.ID = id
}

func (u *User) UpdateName(name string) {
	u.Name = name
}

func (u *User) UpdateEmail(email string) {
	u.Email = email
}

func (u *User) UpdateIsActive(isActive bool) {
	u.IsActive = isActive
}
