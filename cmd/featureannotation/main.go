package main

import (
	"fmt"
	"os"

	facli "github.com/dictyBase/modware-import/internal/featureannotation/cli"
	faclient "github.com/dictyBase/modware-import/internal/featureannotation/client"
	"github.com/dictyBase/modware-import/internal/logger"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "featureannotation",
		Usage: "A command line application for managing feature annotations",
		Flags: []cli.Flag{
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
		Commands: allCommands(),
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func allCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:   "load-feature-annotation",
			Usage:  "Load feature annotations from ArangoDB to feature annotation service",
			Flags:  facli.LoadFeatureAnnotationFlag(),
			Before: faclient.CliSetup,
			Action: facli.LoadFeatureAnnotation,
		},
		// Add the new command definition here
		{
			Name:   "load-csv-to-arangodb",
			Usage:  "Update ArangoDB collection data from CSV file",
			Flags:  facli.LoadCSVToArangodbFlag(),
			Before: faclient.CSVArangodbCliSetup,
			Action: facli.LoadCSVToArangodb,
		},
	}
}
