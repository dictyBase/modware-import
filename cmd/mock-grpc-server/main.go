package main

import (
	"os"

	mockcli "github.com/dictyBase/modware-import/internal/mock-grpc-server/cli"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "mock-grpc-server",
		Usage:    "Mock gRPC servers for integration testing",
		Commands: buildCommands(),
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

// buildCommands creates all CLI commands
func buildCommands() []*cli.Command {
	return []*cli.Command{
		buildAnnotationCommand(),
		buildStockCommand(),
	}
}

// buildAnnotationCommand creates the annotation subcommand
func buildAnnotationCommand() *cli.Command {
	return &cli.Command{
		Name:   "annotation",
		Usage:  "Mock gRPC server for feature annotation service",
		Flags:  buildAnnotationFlags(),
		Action: mockcli.RunServer,
	}
}

// buildStockCommand creates the stock subcommand
func buildStockCommand() *cli.Command {
	return &cli.Command{
		Name:   "stock",
		Usage:  "Mock gRPC server for stock service",
		Flags:  buildStockFlags(),
		Action: mockcli.RunStockServer,
	}
}

// buildAnnotationFlags creates flags for the annotation command
func buildAnnotationFlags() []cli.Flag {
	return []cli.Flag{
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
	}
}

// buildStockFlags creates flags for the stock command
func buildStockFlags() []cli.Flag {
	commonFlags := buildCommonLoggingFlags()
	stockSpecificFlags := buildStockSpecificFlags()
	return append(commonFlags, stockSpecificFlags...)
}

// buildCommonLoggingFlags creates common logging flags
func buildCommonLoggingFlags() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Value:   9560,
			Usage:   "Server port",
			EnvVars: []string{"STOCK_GRPC_PORT"},
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
	}
}

// buildStockSpecificFlags creates stock-specific flags
func buildStockSpecificFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "data-dir",
			Value:   "",
			Usage:   "Pebble data directory (empty for in-memory)",
			EnvVars: []string{"PEBBLE_DATA_DIR"},
		},
		&cli.BoolFlag{
			Name:  "reflection",
			Value: true,
			Usage: "Enable gRPC server reflection",
		},
		&cli.StringFlag{
			Name:  "strain-ontology",
			Value: "dicty_strain_property",
			Usage: "Ontology for strain grouping terms",
		},
		&cli.StringFlag{
			Name:  "strain-term",
			Value: "general strain",
			Usage: "Default ontology term for strains",
		},
		&cli.StringFlag{
			Name:  "plasmid-ontology",
			Value: "plasmid_keywords",
			Usage: "Ontology for plasmid grouping terms",
		},
		&cli.StringFlag{
			Name:  "plasmid-term",
			Value: "cloning vector",
			Usage: "Default ontology term for plasmids",
		},
	}
}
