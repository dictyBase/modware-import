package main

import (
	"log"
	"os"

	"github.com/dictyBase/modware-import/internal/gpadstats"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "gpad-stats",
		Usage: "Generate statistics from Gene Ontology annotation files (GPAD)",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "Path to the GPAD/TSV file",
				Required: true,
			},
		},
		Action: gpadstats.Run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
