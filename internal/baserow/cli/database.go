package cli

import (
	"context"
	"errors"
	"fmt"

	"net/http"

	"github.com/dictyBase/modware-import/internal/baserow/client"
	"github.com/urfave/cli/v2"
)

type accessTokenProperties struct {
	email, password, server string
}

func baserowClient(server string) *client.APIClient {
	conf := client.NewConfiguration()
	conf.Host = server
	conf.Scheme = "https"
	return client.NewAPIClient(conf)
}

func accessToken(args *accessTokenProperties) (string, error) {
	req := baserowClient(
		args.server,
	).UserApi.TokenAuth(
		context.Background(),
	)
	resp, r, err := req.TokenObtainPairWithUser(
		client.TokenObtainPairWithUser{
			Email:    &args.email,
			Password: args.password,
		},
	).Execute()
	defer r.Body.Close()
	if err != nil {
		return "", fmt.Errorf("error in executing API call %s", err)
	}
	if r != nil && r.StatusCode == http.StatusUnauthorized {
		return "", errors.New("unauthrorized access")
	}
	return resp.GetAccessToken(), nil
}

func CreateAccessToken(c *cli.Context) error {
	token, err := accessToken(&accessTokenProperties{
		email:    c.String("email"),
		password: c.String("password"),
		server:   c.String("server"),
	})
	if err != nil {
		return cli.Exit(err, 2)
	}
	fmt.Println(token)
	return nil
}

func CreateAccessTokenFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "email",
			Aliases:  []string{"e"},
			Usage:    "Email of the user",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "password",
			Aliases:  []string{"p"},
			Usage:    "Database password",
			Required: true,
		},
	}
}
