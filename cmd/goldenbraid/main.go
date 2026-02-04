package main

import (
	"log"
	"os"
	"slices"

	"github.com/dictyBase/modware-import/internal/cli/stockcenter"
	loader "github.com/dictyBase/modware-import/internal/load/stockcenter"
	"github.com/urfave/cli/v2"
)

// Common flags shared across all subcommands
var commonFlags = []cli.Flag{
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
				Flags:  slices.Concat(commonFlags, annotationFlags),
			},
			{
				Name:   "plasmid-ontology",
				Usage:  "Associate plasmids with ontology keywords",
				Action: loader.LoadPlasmidOntologyCli,
				Before: stockcenter.SetStockAndS3Clients,
				Flags: slices.Concat(
					commonFlags,
					s3Flags,
					[]cli.Flag{
						&cli.StringFlag{
							Name:  "property",
							Usage: "Property label to filter plasmid keyword rows",
							Value: "keyword",
						},
					},
				),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
