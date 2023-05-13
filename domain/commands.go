package domain

type CreateTeamCommand struct {
	Name        string
	Description string
	OwnerEmail  string
}
