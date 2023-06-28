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
				Name:   "create-access-token",
				Usage:  "Create a new baserow access token",
				Flags:  baserow.CreateAccessTokenFlag(),
				Action: baserow.CreateAccessToken,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
