package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type TeamSrvComponentTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *TeamSrvComponentTestSuite) SetupSuite() {
	conf := Config{
		DbHost:     "localhost",
		DbUser:     "test",
		DbPassword: "test",
		DbName:     "test",
		DbPort:     5432,
		DbSSLMode:  "disable",
	}
	srv := NewServiceWithConfig(conf)
	srv.Setup()

	suite.router = srv.Router
}

func (suite *TeamSrvComponentTestSuite) TearDownSuite() {
	DB.Exec("TRUNCATE teams CASCADE")
}

func (suite *TeamSrvComponentTestSuite) TestCreateTeam() {
	reqBody, _ := json.Marshal(map[string]string{
		"name":        "TestTeam",
		"description": "Team created for test purposes",
		"owner":       "fake@example.com",
	})
	w := httptest.NewRecorder()
	createTeamRequest, _ := http.NewRequest(http.MethodPost, "/teams", bytes.NewBuffer(reqBody))

	suite.router.ServeHTTP(w, createTeamRequest)
	suite.Equal(http.StatusCreated, w.Code)

	var res map[string]string
	err := json.NewDecoder(w.Body).Decode(&res)
	suite.Nil(err)

	w = httptest.NewRecorder()
	fetchTeamRequest, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/teams/%s", res["uuid"]), nil)

	suite.router.ServeHTTP(w, fetchTeamRequest)
	suite.Equal(http.StatusOK, w.Code)
}

func TestTeamSrvComponentTestSuite(t *testing.T) {
	suite.Run(t, new(TeamSrvComponentTestSuite))
}
