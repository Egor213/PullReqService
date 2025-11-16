package integration_test

import (
	"app/internal/entity"
	"net/http"
	"testing"

	. "github.com/Eun/go-hit"
)

func TestCreateTeam_CreateUpdate(t *testing.T) {

	baseTeam := map[string]any{
		"team_name": "backend_test_case",
		"members": []map[string]any{
			{"user_id": "u1000", "username": "Alice", "is_active": true},
			{"user_id": "u1001", "username": "Bob", "is_active": true},
		},
	}

	changedTeam := map[string]any{
		"team_name": baseTeam["team_name"],
		"members": []map[string]any{
			{"user_id": "u1000", "username": "Alice", "is_active": true},
			{"user_id": "u1002", "username": "Bob", "is_active": true},
		},
	}

	testCases := []struct {
		description string
		statusCode  int
		body        map[string]any
	}{
		{
			description: "Create team successfully",
			statusCode:  http.StatusCreated,
			body:        baseTeam,
		},
		{
			description: "Update users the team",
			statusCode:  http.StatusCreated,
			body:        changedTeam,
		},
	}

	for _, tc := range testCases {
		expectedResponse := map[string]any{
			"team": tc.body,
		}
		Test(t,
			Description(tc.description),
			Post(BasePath+"/team/add"),
			Send().Headers("Content-Type").Add("application/json"),
			Send().Body().JSON(tc.body),
			Expect().Status().Equal(int64(tc.statusCode)),
			Expect().Body().JSON().Equal(expectedResponse),
		)
	}
}

func TestCreateTeam_EmptyMembers(t *testing.T) {

	baseTeam := map[string]any{
		"team_name": "backend_test_case",
		"members":   []map[string]any{},
	}
	err := map[string]any{
		"error": map[string]any{
			"code":    "INVALID_REQUEST_PARAMETERS",
			"message": "field Members must be at least 1 characters",
		},
	}
	Test(t,
		Description("Create team empty members"),
		Post(BasePath+"/team/add"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().JSON(baseTeam),
		Expect().Status().Equal(int64(http.StatusBadRequest)),
		Expect().Body().JSON().Equal(err),
	)
}

func TestGetTeam_BaseErrors(t *testing.T) {

	testCases := []struct {
		description      string
		teamName         string
		withAuth         bool
		expectedStatus   IStep
		expectedResponse IStep
	}{
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

	testToken, _ := getAuthToken("u1", string(entity.RoleAdmin), "http://"+Host)

	for _, tc := range testCases {
		authToken := "none"
		if tc.withAuth {
			authToken = testToken
		}
		authHeader := Send().Headers("Authorization").Add("Bearer " + authToken)
		Test(t,
			Description(tc.description),
			Get(BasePath+"/team/get?team_name="+tc.teamName),
			authHeader,
			tc.expectedStatus,
			tc.expectedResponse,
		)
	}
}

func TestGetTeam_OK(t *testing.T) {
	testToken, _ := getAuthToken("u1", string(entity.RoleAdmin), "http://"+Host)
	authHeader := Send().Headers("Authorization").Add("Bearer " + testToken)
	team := map[string]any{
		"team_name": "getTeamTest",
		"members": []map[string]any{
			{"user_id": "getTeamTest1", "username": "getTeamTest1", "is_active": true},
			{"user_id": "getTeamTest2", "username": "getTeamTest2", "is_active": false},
		},
	}
	Test(t,
		Description("OK get team"),
		Get(BasePath+"/team/get?team_name="+"getTeamTest"),
		authHeader,
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().Equal(team),
	)
}
