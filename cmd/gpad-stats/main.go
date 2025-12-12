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
		Commands: []*cli.Command{
			{
				Name:    "gene-count",
				Aliases: []string{"gc"},
				Usage:   "Calculate unique gene counts from GPAD file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Usage:    "Path to the GPAD/TSV file",
						Required: true,
					},
				},
				Action: gpadstats.Run,
			},
			{
				Name:    "gene-count-url",
				Aliases: []string{"gcu"},
				Usage:   "Calculate unique gene counts from GPAD file URL",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "url",
						Aliases:  []string{"u"},
						Usage:    "URL of the gzipped GPAD file",
						Required: true,
					},
				},
				Action: gpadstats.RunURL,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}