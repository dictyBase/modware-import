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
		},
		Commands: []*cli.Command{
			{
				Name:   "create-database-token",
				Usage:  "Create a new baserow database token",
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
