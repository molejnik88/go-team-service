package main

import "github.com/molejnik88/go-team-service/domain"

type CreateTeamRequestBody struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	OwnerEmail  string `json:"owner" binding:"required,email"`
}

type CreateTeamResponseBody struct {
	UUID string `json:"uuid" binding:"required"`
}

type FetchTeamResponseBody struct {
	UUID        string                   `json:"uuid" binding:"required"`
	Name        string                   `json:"name" binding:"required"`
	Description string                   `json:"description"`
	Members     []TeamMemberResponseBody `json:"members" binding:"required"`
}

type TeamMemberResponseBody struct {
	Email   string `json:"email" binding:"required"`
	IsAdmin bool   `json:"is_admin" binding:"required"`
	IsOwner bool   `json:"is_owner" binding:"required"`
}

func (ft *FetchTeamResponseBody) fromTeamModel(team *domain.Team) {
	ft.UUID = team.UUID
	ft.Name = team.Name
	ft.Description = team.Description
	ft.Members = make([]TeamMemberResponseBody, len(team.Members))

	for i, member := range team.Members {
		ft.Members[i] = TeamMemberResponseBody{
			Email:   member.Email,
			IsAdmin: member.IsAdmin,
			IsOwner: member.IsOwner,
		}
	}
}

type AddTeamMemberRequestBody struct {
	UserEmail string `json:"user_email" binding:"required"`
	IsAdmin   bool   `json:"is_admin"`
}
