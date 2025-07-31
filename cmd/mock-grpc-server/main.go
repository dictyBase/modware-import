package main

import (
	"os"

	mockcli "github.com/dictyBase/modware-import/internal/mock-grpc-server/cli"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "mock-grpc-server",
		Usage: "Mock gRPC server for feature annotation service integration testing",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   9000,
				Usage:   "Server port",
				EnvVars: []string{"GRPC_PORT"},
			},
			&cli.StringFlag{
				Name:    "log-level",
				Aliases: []string{"l"},
				Value:   "info",
				Usage:   "Log level (debug, info, warn, error)",
				EnvVars: []string{"LOG_LEVEL"},
			},
			&cli.StringFlag{
				Name:    "log-format",
				Value:   "json",
				Usage:   "Log format (json, text)",
				EnvVars: []string{"LOG_FORMAT"},
			},
			&cli.StringFlag{
				Name:    "log-file",
				Usage:   "Log file path (optional)",
				EnvVars: []string{"LOG_FILE"},
			},
			&cli.StringFlag{
				Name:    "storage-type",
				Aliases: []string{"s"},
				Value:   "leveldb",
				Usage:   "Storage backend type (leveldb uses in-memory storage, memory uses simple map)",
				EnvVars: []string{"STORAGE_TYPE"},
			},
		},
		Action: mockcli.RunServer,
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
