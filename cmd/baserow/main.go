package main

import (
	"fmt"
	"os"

	baserow "github.com/dictyBase/modware-import/internal/baserow/cli"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "baserow",
		Usage: "A command line application for managing baserow instance",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "server",
				Usage:    "address of api server",
				Required: true,
				Aliases:  []string{"s"},
			},
			&cli.StringFlag{
				Name:  "log-level",
				Usage: "Logging level, should be one of debug,warn,info or error",
				Value: "error",
			},
			&cli.StringFlag{
				Name:  "log-format",
				Usage: "Format of log, either of json or text",
				Value: "json",
			},
			&cli.StringFlag{
				Name:  "log-file",
				Usage: "log file for output in addition to stderr",
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "create-access-token",
				Usage:  "Create a new baserow access token",
				Flags:  baserow.CreateAccessTokenFlag(),
				Action: baserow.CreateAccessToken,
			},
			{
				Name:   "create-database-token",
				Usage:  "Create a baserow database token",
				Flags:  baserow.CreateDatabaseTokenFlag(),
				Action: baserow.CreateDatabaseToken,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
