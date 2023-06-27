package cli

import (
	"github.com/urfave/cli/v2"
)

func CreateDatabaseTokenFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "user",
			Aliases:  []string{"u"},
			Usage:    "Database user",
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
