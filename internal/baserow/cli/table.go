package cli

import (
	"context"
	"fmt"

	"github.com/dictyBase/modware-import/internal/baserow/client"
	"github.com/urfave/cli/v2"
)

func CreateTable(c *cli.Context) error {
	bclient := baserowClient(c.String("server"))
	authCtx := context.WithValue(
		context.Background(),
		client.ContextAccessToken,
		c.String("token"),
	)
	_, resp, err := bclient.DatabaseTablesApi.CreateDatabaseTable(authCtx, int32(c.Int("database-id"))).
		Execute()
	defer resp.Body.Close()
	if err != nil {
		return cli.Exit(
			fmt.Errorf("error in creating table %s", err), 2,
		)
	}
	return nil
}

func CreateTableFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "token",
			Aliases:  []string{"t"},
			Usage:    "database token with write privilege",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "database-id",
			Usage:    "Database id",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "table",
			Aliases:  []string{"t"},
			Usage:    "Database table",
			Required: true,
		},
	}
}
