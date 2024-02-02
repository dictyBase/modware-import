package main

import (
	"fmt"
	"os"

	"github.com/dictyBase/modware-import/internal/logger"
	logto "github.com/dictyBase/modware-import/internal/logto/cli"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/jellydator/ttlcache/v3"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "logto",
		Usage:  "A command line application for logto instance management",
		Before: logger.SetupCliLogger,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "endpoint",
				Usage:    "endpoint of api server",
				Required: true,
				Aliases:  []string{"e"},
			},
			&cli.StringFlag{
				Name:     "app-secret",
				Usage:    "logto application secret to access the api",
				Aliases:  []string{"s"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "app-id",
				Usage:    "logto application identifier to access the api",
				Aliases:  []string{"a"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "api-resource",
				Usage:    "logto api resource name",
				Aliases:  []string{"r"},
				Required: true,
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
				Usage:  "Create a new logto access token",
				Action: logto.CreateAccessToken,
			},
			{
				Name:   "import-user",
				Usage:  "Import user from an input file",
				Action: logto.ImportUser,
				Before: setupTTLCache,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "input",
						Usage:    "input csv file with users information",
						Aliases:  []string{"i"},
						Required: true,
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func setupTTLCache(cltx *cli.Context) error {
	tcache := ttlcache.New[string, string]()
	registry.SetTTLCache(tcache)
	reader, err := os.Open(cltx.String("input"))
	if err != nil {
		return fmt.Errorf("error in reading file %s", err)
	}
	registry.SetReader("USER_INPUT", reader)
	return nil
}
