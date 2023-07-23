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

func (suite *TeamSrvComponentTestSuite) TearDownTest() {
	DB.Exec("TRUNCATE teams CASCADE")
}

func (suite *TeamSrvComponentTestSuite) createTeam() string {
	reqBody, _ := json.Marshal(map[string]string{
		"name":        "TestTeam",
		"description": "Team created for test purposes",
		"owner":       "fake@example.com",
	})
	w := httptest.NewRecorder()
	createTeamRequest, _ := http.NewRequest(http.MethodPost, "/teams", bytes.NewReader(reqBody))

	suite.router.ServeHTTP(w, createTeamRequest)
	suite.Equal(http.StatusCreated, w.Code)

	var res map[string]string
	err := json.NewDecoder(w.Body).Decode(&res)
	suite.Nil(err)

	return res["uuid"]
}

func (suite *TeamSrvComponentTestSuite) TestCreateTeam() {
	teamUUID := suite.createTeam()

	w := httptest.NewRecorder()
	fetchTeamRequest, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/teams/%s", teamUUID), nil)

	suite.router.ServeHTTP(w, fetchTeamRequest)
	suite.Equal(http.StatusOK, w.Code)
}

func (suite *TeamSrvComponentTestSuite) TestAddMember() {
	var w *httptest.ResponseRecorder
	teamUUID := suite.createTeam()

	reqBody, _ := json.Marshal(map[string]string{
		"user_email": "fake_mamber@example.com",
	})
	w = httptest.NewRecorder()
	addMember, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/teams/%s/add_member", teamUUID), bytes.NewReader(reqBody))
	suite.router.ServeHTTP(w, addMember)

	suite.Require().Equal(http.StatusCreated, w.Code)

	w = httptest.NewRecorder()
	fetchTeam, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/teams/%s", teamUUID), nil)
	suite.router.ServeHTTP(w, fetchTeam)
	suite.Equal(http.StatusOK, w.Code)

	var res FetchTeamResponseBody
	err := json.NewDecoder(w.Body).Decode(&res)
	suite.Nil(err)
	suite.Len(res.Members, 2)
}

func TestTeamSrvComponentTestSuite(t *testing.T) {
	suite.Run(t, new(TeamSrvComponentTestSuite))
}
