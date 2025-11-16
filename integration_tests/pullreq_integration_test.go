package integration_test

import (
	"app/internal/entity"
	"net/http"
	"testing"

	. "github.com/Eun/go-hit"
)

func TestPullReq_BaseErrors(t *testing.T) {
	req := map[string]any{
		"pull_request_id":   "create_test_pr",
		"pull_request_name": "test",
		"author_id":         "dfgfffffffffff",
	}
	reqReassign := map[string]any{
		"pull_request_id": "fffffffff",
		"old_reviewer_id": "fffffffffffffffffff",
	}
	testCases := []struct {
		description      string
		body             map[string]any
		bodyReassing     map[string]any
		withAuth         bool
		expectedStatus   IStep
		expectedResponse IStep
		isAdmin          bool
	}{
		{
			description:      "User not found",
			body:             req,
			bodyReassing:     reqReassign,
			withAuth:         true,
			isAdmin:          true,
			expectedStatus:   Expect().Status().Equal(http.StatusNotFound),
			expectedResponse: Expect().Body().JSON().JQ(".error.code").Equal("NOT_FOUND"),
		},
		{
			description:      "Unauthorized user",
			body:             map[string]any{},
			bodyReassing:     map[string]any{},
			withAuth:         false,
			isAdmin:          false,
			expectedStatus:   Expect().Status().Equal(http.StatusUnauthorized),
			expectedResponse: Expect().Body().JSON().JQ(".error.code").Equal("NOT_FOUND"),
		},
		{
			description:      "Invalid params",
			body:             map[string]any{},
			bodyReassing:     map[string]any{},
			withAuth:         true,
			isAdmin:          true,
			expectedStatus:   Expect().Status().Equal(http.StatusBadRequest),
			expectedResponse: Expect().Body().JSON().JQ(".error.code").Equal("INVALID_REQUEST_PARAMETERS"),
		},
		{
			description:      "Forbidden",
			body:             map[string]any{},
			bodyReassing:     map[string]any{},
			withAuth:         true,
			isAdmin:          false,
			expectedStatus:   Expect().Status().Equal(http.StatusForbidden),
			expectedResponse: Expect().Body().JSON().JQ(".error.code").Equal("FORBIDDEN"),
		},
	}

	testToken, _ := getAuthToken("u1", string(entity.RoleAdmin), "http://"+Host)
	testTokenUser, _ := getAuthToken("u1", string(entity.RoleUser), "http://"+Host)

	for _, tc := range testCases {
		authToken := "none"
		if tc.withAuth {
			if tc.isAdmin {
				authToken = testToken
			} else {
				authToken = testTokenUser
			}
		}

		authHeader := Send().Headers("Authorization").Add("Bearer " + authToken)
		Test(t,
			Description(tc.description+" create"),
			Post(BasePath+"/pullRequest/create"),
			authHeader,
			Send().Headers("Content-Type").Add("application/json"),
			Send().Body().JSON(tc.body),
			tc.expectedStatus,
			tc.expectedResponse,
		)
		Test(t,
			Description(tc.description+" reassign"),
			Post(BasePath+"/pullRequest/reassign"),
			authHeader,
			Send().Headers("Content-Type").Add("application/json"),
			Send().Body().JSON(tc.bodyReassing),
			tc.expectedStatus,
			tc.expectedResponse,
		)
		if tc.description != "User not found" {
			Test(t,
				Description(tc.description+" merge"),
				Post(BasePath+"/pullRequest/merge"),
				authHeader,
				Send().Headers("Content-Type").Add("application/json"),
				Send().Body().JSON(tc.body),
				tc.expectedStatus,
				tc.expectedResponse,
			)
		}
	}
}

func TestCreatePR_OKAndExists(t *testing.T) {
	req := map[string]any{
		"pull_request_id":   "create_test_pr",
		"pull_request_name": "test",
		"author_id":         "TestCreatePRUser",
	}
	resp := map[string]any{
		"pr": map[string]any{
			"pull_request_id":   "create_test_pr",
			"pull_request_name": "test",
			"author_id":         "TestCreatePRUser",
			"status":            "OPEN",
			"assigned_reviewers": []string{
				"TestCreatePRUserReviewer1",
			},
		},
	}

	testCases := []struct {
		description      string
		body             map[string]any
		expectedStatus   IStep
		expectedResponse IStep
	}{
		{
			description:      "OK create pr",
			body:             req,
			expectedStatus:   Expect().Status().Equal(http.StatusCreated),
			expectedResponse: Expect().Body().JSON().Equal(resp),
		},
		{
			description:      "PR already exists",
			body:             req,
			expectedStatus:   Expect().Status().Equal(http.StatusConflict),
			expectedResponse: Expect().Body().JSON().JQ(".error.code").Equal("PR_EXISTS"),
		},
	}

	testToken, _ := getAuthToken("u1", string(entity.RoleAdmin), "http://"+Host)

	for _, tc := range testCases {
		authHeader := Send().Headers("Authorization").Add("Bearer " + testToken)
		Test(t,
			Description(tc.description),
			Post(BasePath+"/pullRequest/create"),
			Send().Headers("Content-Type").Add("application/json"),
			Send().Body().JSON(tc.body),
			authHeader,
			tc.expectedStatus,
			tc.expectedResponse,
		)
	}
}

func TestMergePR_OKAndExists(t *testing.T) {
	req := map[string]any{
		"pull_request_id": "TestMergePR",
	}

	testCases := []struct {
		description      string
		body             map[string]any
		expectedStatus   IStep
		expectedResponse IStep
	}{
		{
			description:      "OK create pr",
			body:             req,
			expectedStatus:   Expect().Status().Equal(http.StatusOK),
			expectedResponse: Expect().Body().JSON().JQ(".pr.status").Equal("MERGED"),
		},
		{
			description:      "Idempotenty",
			body:             req,
			expectedStatus:   Expect().Status().Equal(http.StatusOK),
			expectedResponse: Expect().Body().JSON().JQ(".pr.status").Equal("MERGED"),
		},
	}

	testToken, _ := getAuthToken("u1", string(entity.RoleAdmin), "http://"+Host)

	for _, tc := range testCases {
		authHeader := Send().Headers("Authorization").Add("Bearer " + testToken)
		Test(t,
			Description(tc.description),
			Post(BasePath+"/pullRequest/merge"),
			Send().Headers("Content-Type").Add("application/json"),
			Send().Body().JSON(tc.body),
			authHeader,
			tc.expectedStatus,
			tc.expectedResponse,
		)
	}
}

func TestReassignUser_OKAndMerged(t *testing.T) {
	reqReassign := map[string]any{
		"pull_request_id": "TestReassignPR1",
		"old_reviewer_id": "TestReassignUser2",
	}
	testCases := []struct {
		description      string
		body             map[string]any
		expectedStatus   IStep
		expectedResponse IStep
	}{
		{
			description:      "OK reassign reviewer",
			body:             reqReassign,
			expectedStatus:   Expect().Status().Equal(http.StatusOK),
			expectedResponse: Expect().Body().JSON().JQ(".replaced_by").Equal("TestReassignUser4"),
		},
		{
			description:      "reviewer is not assigned to this PR",
			body:             reqReassign,
			expectedStatus:   Expect().Status().Equal(http.StatusNotFound),
			expectedResponse: Expect().Body().JSON().JQ(".error.code").Equal("NOT_ASSIGNED"),
		},
		{
			description: "merged pr",
			body: map[string]any{
				"pull_request_id": "TestReassignPR2",
				"old_reviewer_id": "TestReassignUser2",
			},
			expectedStatus:   Expect().Status().Equal(http.StatusConflict),
			expectedResponse: Expect().Body().JSON().JQ(".error.code").Equal("PR_MERGED"),
		},
	}

	testToken, _ := getAuthToken("u1", string(entity.RoleAdmin), "http://"+Host)

	for _, tc := range testCases {
		authHeader := Send().Headers("Authorization").Add("Bearer " + testToken)
		Test(t,
			Description(tc.description),
			Post(BasePath+"/pullRequest/reassign"),
			authHeader,
			Send().Headers("Content-Type").Add("application/json"),
			Send().Body().JSON(tc.body),
			tc.expectedStatus,
			tc.expectedResponse,
		)
	}

}
