package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/molejnik88/go-team-service/adapters"
	"github.com/molejnik88/go-team-service/domain"
	"github.com/molejnik88/go-team-service/service_layer"
)

func createTeam(c *gin.Context) {
	var createTeam CreateTeamRequestBody
	uow := adapters.NewSqlUnitOfWork(DB)

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

	uuid, err := service_layer.CreateTeam(command, uow)

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
	team := new(domain.Team)
	result := DB.Preload("Members").First(team, "uuid = ?", uuid)
	if err := result.Error; err != nil {
		// TODO: Not found error?
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	responseBody.fromTeamModel(team)
	c.JSON(http.StatusOK, responseBody)
}
