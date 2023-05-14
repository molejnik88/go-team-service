package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTeam(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		reqBody, _ := json.Marshal(map[string]string{
			"name":        "TestTeam",
			"description": "Team created for test purposes",
			"owner":       "fake@example.com",
		})
		router := setupRouter()
		w := httptest.NewRecorder()
		createTeamRequest, _ := http.NewRequest(http.MethodPost, "/teams", bytes.NewBuffer(reqBody))

		router.ServeHTTP(w, createTeamRequest)
		assert.Equal(t, http.StatusCreated, w.Code)

		var res map[string]string
		err := json.NewDecoder(w.Body).Decode(&res)
		assert.Nil(t, err)

		w = httptest.NewRecorder()
		fetchTeamRequest, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/teams/%s", res["uuid"]), nil)

		router.ServeHTTP(w, fetchTeamRequest)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
