package domain

type CreateTeamCommand struct {
	Name        string
	Description string
	OwnerEmail  string
}

type AddTeamMember struct {
	TeamUUID string
	Email    string
	IsAdmin  bool
}
