package integration_test

import (
	"net/http"

	. "github.com/Eun/go-hit"
)

var (
	Host     = "app:8080"
	BasePath = "http://" + Host + "/api/v1"
)

func getAuthToken(username, role, url string) (string, error) {
	var token string
	var err error

	body := map[string]any{
		"user_id": username,
		"role":    role,
	}
	err = Do(
		Post(url+"/auth/login"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().JSON(body),
		Expect().Status().Equal(http.StatusOK),
		Store().Response().Body().JSON().JQ(".access_token").In(&token),
	)
	if err != nil {
		return "", err
	}

	return token, nil
}
