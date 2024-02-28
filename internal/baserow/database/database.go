package database

import (
	"context"
	"errors"
	"fmt"

	"net/http"

	"github.com/dictyBase/modware-import/internal/baserow/client"
)

type AccessTokenProperties struct {
	Email, Password, Server string
}

func BaserowClient(server string) *client.APIClient {
	conf := client.NewConfiguration()
	conf.Host = server
	conf.Scheme = "https"
	return client.NewAPIClient(conf)
}

func AccessToken(
	args *AccessTokenProperties,
) (*client.CreateUser200Response, error) {
	req := BaserowClient(
		args.Server,
	).UserApi.TokenAuth(
		context.Background(),
	)
	resp, r, err := req.TokenObtainPairWithUser(
		client.TokenObtainPairWithUser{
			Email:    &args.Email,
			Password: args.Password,
		},
	).Execute()
	defer r.Body.Close()
	if err != nil {
		return resp, fmt.Errorf("error in executing API call %s", err)
	}
	if r != nil && r.StatusCode == http.StatusUnauthorized {
		return resp, errors.New("unauthrorized access")
	}
	return resp, nil
}
