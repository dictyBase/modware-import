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
		Name:  "logto",
		Usage: "A command line application for logto instance management",
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
		Before: func(c *cli.Context) error {
			l, err := logger.NewCliLogger(c)
			if err != nil {
				return fmt.Errorf("error in getting a new logger %s", err)
			}
			registry.SetLogger(l)
			tcache := ttlcache.New[string, int]()
			registry.SetTTLCache(tcache)
			return nil
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
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
