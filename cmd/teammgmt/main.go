package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/molejnik88/go-team-service/adapters"
	"github.com/molejnik88/go-team-service/domain"
	"github.com/molejnik88/go-team-service/service_layer"
)

var uow adapters.InMemoryUOW = adapters.InMemoryUOW{
	Repository: &adapters.InMemoryRepository{
		Teams: make(map[string]domain.Team),
	},
}

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

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/teams", createTeam)
	r.GET("/teams/:uuid", fetchTeam)

	return r
}

func main() {
	r := setupRouter()
	r.Run() // listen and serve on 0.0.0.0:8080
}

func createTeam(c *gin.Context) {
	var createTeam CreateTeamRequestBody
	// TODO: custom error messages for validators
	if err := c.ShouldBindJSON(&createTeam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	command := &domain.CreateTeamCommand{
		Name:        createTeam.Name,
		Description: createTeam.Description,
		OwnerEmail:  createTeam.OwnerEmail,
	}

	uuid, err := service_layer.CreateTeam(command, &uow)

	// TODO: more granular errors
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateTeamResponseBody{UUID: uuid})
}

func fetchTeam(c *gin.Context) {
	responseBody := &FetchTeamResponseBody{}
	uuid := c.Param("uuid")

	// TODO: create a query, don't use repository
	team, err := uow.Teams().Get(uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	responseBody.fromTeamModel(team)
	c.JSON(http.StatusOK, responseBody)
}
