package integration_test

import (
	"app/internal/entity"
	"net/http"
	"testing"

	. "github.com/Eun/go-hit"
)

func TestSetIsActive_SetValues(t *testing.T) {
	testToken, _ := getAuthToken("u1", string(entity.RoleAdmin), "http://"+Host)
	authHeader := Send().Headers("Authorization").Add("Bearer " + testToken)

	testCases := []struct {
		name     string
		isActive bool
	}{
		{name: "Set active true", isActive: true},
		{name: "Set active false", isActive: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := map[string]any{
				"user_id":   "userForSetActive",
				"is_active": tc.isActive,
			}

			resp := map[string]any{
				"user": map[string]any{
					"user_id":   "userForSetActive",
					"username":  "userForSetActive",
					"team_name": "testing",
					"is_active": tc.isActive,
				},
			}

			Test(t,
				Description(tc.name),
				Post(BasePath+"/users/setIsActive"),
				authHeader,
				Send().Headers("Content-Type").Add("application/json"),
				Send().Body().JSON(req),
				Expect().Status().Equal(http.StatusOK),
				Expect().Body().JSON().Equal(resp),
			)
		})
	}
}

func TestUsers_BaseErrors(t *testing.T) {
	testCases := []struct {
		description      string
		body             map[string]any
		reviewParam      string
		withAuth         bool
		expectedStatus   IStep
		expectedResponse IStep
		isAdmin          bool
	}{
		{
			description: "User not found",
			body: map[string]any{
				"user_id":   "dfagfshgjkljhgfdsa",
				"is_active": true,
			},
			reviewParam:      "fffffffffffffffffff",
			withAuth:         true,
			isAdmin:          true,
			expectedStatus:   Expect().Status().Equal(http.StatusNotFound),
			expectedResponse: Expect().Body().JSON().JQ(".error.code").Equal("NOT_FOUND"),
		},
		{
			description: "Unauthorized user",
			body: map[string]any{
				"user_id":   "",
				"is_active": true,
			},
			reviewParam:      "TestGetReview",
			withAuth:         false,
			isAdmin:          false,
			expectedStatus:   Expect().Status().Equal(http.StatusUnauthorized),
			expectedResponse: Expect().Body().JSON().JQ(".error.code").Equal("NOT_FOUND"),
		},
		{
			description:      "Invalid params",
			body:             map[string]any{},
			reviewParam:      "",
			withAuth:         true,
			isAdmin:          true,
			expectedStatus:   Expect().Status().Equal(http.StatusBadRequest),
			expectedResponse: Expect().Body().JSON().JQ(".error.code").Equal("INVALID_REQUEST_PARAMETERS"),
		},
		{
			description:      "Forbidden",
			body:             map[string]any{},
			reviewParam:      "TestGetReview",
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
			Description(tc.description+" setIsActive"),
			Post(BasePath+"/users/setIsActive"),
			authHeader,
			Send().Headers("Content-Type").Add("application/json"),
			Send().Body().JSON(tc.body),
			tc.expectedStatus,
			tc.expectedResponse,
		)
		if tc.description != "Forbidden" {
			Test(t,
				Description(tc.description+" getReview"),
				Get(BasePath+"/users/getReview?user_id="+tc.reviewParam),
				authHeader,
				tc.expectedStatus,
				tc.expectedResponse,
			)
		}
	}
}

func TestGetReview_OK(t *testing.T) {
	testToken, _ := getAuthToken("u1", string(entity.RoleAdmin), "http://"+Host)
	authHeader := Send().Headers("Authorization").Add("Bearer " + testToken)

	Test(t,
		Description("getReview"),
		Get(BasePath+"/users/getReview?user_id=TestGetReview"),
		authHeader,
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().JQ(".pull_requests").Len().Equal(2),
	)
}
