package cli

import (
	"context"
	"fmt"

	"net/http"

	"github.com/dictyBase/modware-import/internal/baserow/client"
	"github.com/urfave/cli/v2"
)

func CreateDatabaseToken(c *cli.Context) error {
	conf := client.NewConfiguration()
	conf.Host = c.String("server")
	conf.Scheme = "https"
	bclient := client.NewAPIClient(conf)
	req := bclient.UserApi.TokenAuth(context.Background())
	resp, r, err := req.TokenObtainPairWithUser(
		client.TokenObtainPairWithUser{
			Email:    client.PtrString(c.String("email")),
			Password: c.String("password"),
		},
	).Execute()
	defer r.Body.Close()
	if err != nil {
		return cli.Exit(
			fmt.Sprintf("error in executing API call %s", err),
			2,
		)
	}
	if r != nil && r.StatusCode == http.StatusUnauthorized {
		return cli.Exit("unauthrorized access", 2)
	}
	fmt.Printf("access token %s\n", resp.GetAccessToken())
	return nil
}

func CreateDatabaseTokenFlag() []cli.Flag {
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
