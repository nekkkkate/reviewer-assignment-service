package models

import "errors"

type Team struct {
	ID      int                 `json:"id"`
	Name    string              `json:"name"`
	Members map[int]*TeamMember `json:"members"`
}

func NewTeam(name string) *Team {
	return &Team{
		Name:    name,
		Members: make(map[int]*TeamMember),
	}
}
func (t *Team) AddMember(member *TeamMember) error {
	if t.IsMemberInTeam(member.UserID) {
		return ErrMemberAlreadyInTeam
	}
	t.Members[member.UserID] = member
	return nil
}

func (t *Team) RemoveMember(member *TeamMember) error {
	if !t.IsMemberInTeam(member.UserID) {
		return ErrMemberNotInTeam
	}
	delete(t.Members, member.UserID)
	return nil
}

func (t *Team) IsMemberInTeam(memberID int) bool {
	_, exists := t.Members[memberID]
	return exists
}

func (t *Team) GetMembers() []*TeamMember {
	var members []*TeamMember
	for _, member := range t.Members {
		members = append(members, member)
	}
	return members
}

func (t *Team) IsEmpty() bool {
	return len(t.Members) == 0
}

func (t *Team) GetMemberCount() int {
	return len(t.Members)
}

func (t *Team) SetId(id int) {
	t.ID = id
}

func (t *Team) UpdateName(name string) {
	t.Name = name
}

func (t *Team) GetActiveMembers() []*TeamMember {
	var activeMembers []*TeamMember
	for _, member := range t.Members {
		if member.IsActive {
			activeMembers = append(activeMembers, member)
		}
	}
	return activeMembers
}

var (
	ErrMemberAlreadyInTeam = errors.New("member already in team")
	ErrMemberNotInTeam     = errors.New("member not in team")
)
