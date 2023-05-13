package domain

type Team struct {
	UUID        string
	Name        string
	Description string
	Members     []TeamMember
}

type TeamMember struct {
	Email   string
	IsAdmin bool
	IsOwner bool
}
