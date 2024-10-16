package main

import (
	"fmt"
	"os"

	contentcli "github.com/dictyBase/modware-import/internal/content/cli"
	"github.com/dictyBase/modware-import/internal/content/client"
	"github.com/dictyBase/modware-import/internal/logger"
	uniprotcli "github.com/dictyBase/modware-import/internal/uniprot/cli"
	uniprotclient "github.com/dictyBase/modware-import/internal/uniprot/client"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "loader",
		Usage:  "A command line application for loading content data and Uniprot mappings",
		Before: logger.SetupCliLogger,
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
		Commands: []*cli.Command{
			{
				Name:   "content-data",
				Usage:  "load content data from s3 storage",
				Action: contentcli.LoadContent,
				Flags:  contentcli.ContentLoaderFlags(),
				Before: client.CliSetup,
			},
			{
				Name:   "uniprot-mappings",
				Usage:  "load Uniprot mappings",
				Action: uniprotcli.LoadUniprotMappings,
				Flags:  uniprotcli.UniprotFlags(),
				Before: uniprotclient.SetRedisClient,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
