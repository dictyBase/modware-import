package main

import (
	"log"
	"os"
	"slices"

	"github.com/dictyBase/modware-import/internal/cli/stockcenter"
	loader "github.com/dictyBase/modware-import/internal/load/stockcenter"
	"github.com/urfave/cli/v2"
)

// Stock connection flags (shared by both commands)
var stockConnectionFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "stock-grpc-host",
		Usage:   "gRPC host address for stock service",
		Value:   "stock-api",
		EnvVars: []string{"STOCK_API_SERVICE_HOST"},
	},
	&cli.StringFlag{
		Name:    "stock-grpc-port",
		Usage:   "gRPC port for stock service",
		Value:   "9560",
		EnvVars: []string{"STOCK_API_SERVICE_PORT"},
	},
}

// Input flags (for file-based commands like inventory)
var inputFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     "input",
		Usage:    "Input file path (local path or filename in bucket)",
		Required: true,
	},
	&cli.StringFlag{
		Name:  "input-source",
		Usage: "Source of input file (folder or bucket)",
		Value: "bucket",
	},
	&cli.StringFlag{
		Name:  "s3-bucket",
		Usage: "S3 bucket for input files",
		Value: "dictybase",
	},
	&cli.StringFlag{
		Name:  "s3-bucket-path",
		Usage: "Path inside S3 bucket for input files",
		Value: "import/data/stockcenter",
	},
}

// Ontology-specific flags (for plasmid-ontology command)
var ontologyFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "ontology-term",
		Usage: "Target ontology term to assign to plasmids",
		Value: "vector",
	},
	&cli.IntFlag{
		Name:  "batch-size",
		Usage: "Number of plasmids to fetch per API call",
		Value: 100,
	},
}

// S3 configuration flags
var s3Flags = []cli.Flag{
	&cli.StringFlag{
		Name:    "s3-server",
		Usage:   "S3 server endpoint",
		Value:   "minio",
		EnvVars: []string{"MINIO_SERVICE_HOST"},
	},
	&cli.StringFlag{
		Name:    "s3-server-port",
		Usage:   "S3 server port",
		EnvVars: []string{"MINIO_SERVICE_PORT"},
	},
	&cli.StringFlag{
		Name:    "access-key",
		Usage:   "S3 access key",
		EnvVars: []string{"user"},
	},
	&cli.StringFlag{
		Name:    "secret-key",
		Usage:   "S3 secret key",
		EnvVars: []string{"pass"},
	},
}

// Annotation service flags
var annotationFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "annotation-grpc-host",
		Usage: "gRPC host address for annotation service",
		Value: "annotation-api",
	},
	&cli.StringFlag{
		Name:  "annotation-grpc-port",
		Usage: "gRPC port for annotation service",
		Value: "9560",
	},
}

// Logging flags
var loggingFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "log-level",
		Usage: "Log level (debug, info, warn, error)",
		Value: "info",
	},
	&cli.StringFlag{
		Name:  "log-format",
		Usage: "Log format (json, text)",
		Value: "json",
	},
}

func main() {
	app := &cli.App{
		Name:  "goldenbraid",
		Usage: "GoldenBraid data loading tools",
		Commands: []*cli.Command{
			{
				Name:   "inventory",
				Usage:  "Load GoldenBraid inventory",
				Action: stockcenter.LoadGoldenBraidInventory,
				Before: stockcenter.SetClients,
				Flags: slices.Concat(
					inputFlags,
					s3Flags,
					stockConnectionFlags,
					annotationFlags,
					loggingFlags,
				),
			},
			{
				Name:   "plasmid-ontology",
				Usage:  "Update GB vector plasmids with ontology term",
				Action: loader.LoadPlasmidOntologyCli,
				Before: stockcenter.SetStockClientWrapper,
				Flags: slices.Concat(
					stockConnectionFlags,
					ontologyFlags,
					loggingFlags,
				),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
