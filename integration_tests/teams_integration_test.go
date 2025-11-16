package integration_test

import (
	"app/internal/entity"
	"net/http"
	"testing"

	. "github.com/Eun/go-hit"
)

const (
	host     = "app:8080"
	basePath = "http://" + host + "/api/v1"
)

func TestCreateTeam(t *testing.T) {
	testCases := []struct {
		description      string
		teamName         string
		body             map[string]any
		expectedStatus   IStep
		expectedResponse IStep
	}{
		{
			description: "Create team successfully",
			teamName:    "backend",
			body: map[string]any{
				"team_name": "backend",
				"members": []map[string]any{
					{
						"user_id":   "u1",
						"username":  "Alice",
						"is_active": true,
					},
					{
						"user_id":   "u2",
						"username":  "Bob",
						"is_active": true,
					},
				},
			},
			expectedStatus:   Expect().Status().Equal(http.StatusCreated),
			expectedResponse: Expect().Body().JSON().JQ(".team.team_name").Equal("backend"),
		},
	}

	for _, tc := range testCases {
		Test(t,
			Description(tc.description),
			Post(basePath+"/team/add"),
			Send().Headers("Content-Type").Add("application/json"),
			Send().Body().JSON(tc.body),
			tc.expectedStatus,
			tc.expectedResponse,
		)
	}
}

func TestGetTeam(t *testing.T) {
	testCases := []struct {
		description      string
		teamName         string
		withAuth         bool
		expectedStatus   IStep
		expectedResponse IStep
	}{
		{
			description:      "Get existing team",
			teamName:         "backend",
			withAuth:         true,
			expectedStatus:   Expect().Status().Equal(http.StatusOK),
			expectedResponse: Expect().Body().JSON().JQ(".team_name").Equal("backend"),
		},
		{
			description:      "Get non-existing team",
			teamName:         "nonexistent",
			withAuth:         true,
			expectedStatus:   Expect().Status().Equal(http.StatusNotFound),
			expectedResponse: Expect().Body().JSON().JQ(".error.code").Equal("NOT_FOUND"),
		},
		{
			description:      "Get team without auth",
			teamName:         "backend",
			withAuth:         false,
			expectedStatus:   Expect().Status().Equal(http.StatusUnauthorized),
			expectedResponse: Expect().Body().JSON().JQ(".error").Len().GreaterThan(0),
		},
	}

	testToken, _ := getAuthToken("u1", string(entity.RoleAdmin), "http://"+host)

	for _, tc := range testCases {
		authToken := "none"
		if tc.withAuth {
			authToken = testToken
		}
		authHeader := Send().Headers("Authorization").Add("Bearer " + authToken)
		Test(t,
			Description(tc.description),
			Get(basePath+"/team/get?team_name="+tc.teamName),
			authHeader,
			tc.expectedStatus,
			tc.expectedResponse,
		)
	}
}
