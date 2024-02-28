package main

import (
	"fmt"
	"os"

	baserow "github.com/dictyBase/modware-import/internal/baserow/cli"
	"github.com/dictyBase/modware-import/internal/logger"
	"github.com/dictyBase/modware-import/internal/registry"
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
		Before: func(c *cli.Context) error {
			l, err := logger.NewCliLogger(c)
			if err != nil {
				return fmt.Errorf("error in getting a new logger %s", err)
			}
			registry.SetLogger(l)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:   "create-access-token",
				Usage:  "Create a new baserow access token",
				Flags:  baserow.CreateAccessTokenFlag(),
				Action: baserow.CreateAccessToken,
			},
			{
				Name:   "create-phenotype-table",
				Usage:  "Create a baserow table with phenotype fields preset",
				Flags:  baserow.CreatePhenotypeTableFlag(),
				Action: baserow.CreatePhenoTableHandler,
			},
			{
				Name:   "create-ontology-table",
				Usage:  "Create a baserow table with ontology fields preset",
				Flags:  baserow.CreateOntologyTableFlag(),
				Action: baserow.CreateOntologyTableHandler,
			},
			{
				Name:   "load-ontology",
				Usage:  "load ontology in a baserow table",
				Flags:  baserow.LoadOntologyToTableFlag(),
				Action: baserow.LoadOntologyToTable,
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
