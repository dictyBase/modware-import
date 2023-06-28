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

func CreateDatabaseToken(c *cli.Context) error {
	atoken, err := accessToken(&accessTokenProperties{
		email:    c.String("email"),
		password: c.String("password"),
		server:   c.String("server"),
	})
	if err != nil {
		return cli.Exit(fmt.Errorf("error in creating access token %s", err), 2)
	}
	bclient := baserowClient(c.String("server"))
	authCtx := context.WithValue(
		context.Background(),
		client.ContextAccessToken,
		atoken,
	)
	wlist, r, err := bclient.WorkspacesApi.ListWorkspaces(authCtx).
		Execute()
	defer r.Body.Close()
	if err != nil {
		return cli.Exit(
			fmt.Errorf("error in executing list workspaces API call %s", err),
			2,
		)
	}
	for _, v := range wlist {
		fmt.Println(v.GetName())
	}
	return nil
}

func CreateDatabaseTokenFlag() []cli.Flag {
	aflags := CreateAccessTokenFlag()
	return append(aflags, []cli.Flag{
		&cli.StringFlag{
			Name:     "workspace",
			Aliases:  []string{"w"},
			Usage:    "Only tables under this workspaces can be accessed",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "token name",
			Required: true,
		},
	}...)
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
