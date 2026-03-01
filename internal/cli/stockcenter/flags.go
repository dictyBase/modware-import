package stockcenter

import (
	"github.com/dictyBase/modware-import/internal/config"
	"github.com/urfave/cli/v2"
)

// StockConnectionFlags returns flags for stock service connection
func StockConnectionFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "stock-grpc-host",
			Usage:   "gRPC host address for stock service",
			EnvVars: []string{"STOCK_API_SERVICE_HOST"},
		},
		&cli.StringFlag{
			Name:    "stock-grpc-port",
			Usage:   "gRPC port for stock service",
			EnvVars: []string{"STOCK_API_SERVICE_PORT"},
		},
	}
}

// InputFlags returns flags for input file configuration
func InputFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "input",
			Usage: "Input file path (local path or filename in bucket)",
			Value: "goldenbraid.csv",
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
}

// OntologyFlags returns flags for ontology configuration
func OntologyFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "ontology-term",
			Usage: "Target ontology term to assign to plasmids",
			Value: "vector",
		},
		&cli.IntFlag{
			Name:  "batch-size",
			Usage: "Number of plasmids to fetch per API call",
			Value: config.DefaultPlasmidBatchSize,
		},
	}
}

// S3Flags returns flags for S3 server configuration
func S3Flags() []cli.Flag {
	return []cli.Flag{
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
}

// AnnotationFlags returns flags for annotation service connection
func AnnotationFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "annotation-grpc-host",
			Usage:   "gRPC host address for annotation service",
			Value:   "annotation-api",
			EnvVars: []string{"ANNOTATION_API_SERVICE_HOST"},
		},
		&cli.StringFlag{
			Name:    "annotation-grpc-port",
			Usage:   "gRPC port for annotation service",
			EnvVars: []string{"ANNOTATION_API_SERVICE_PORT"},
		},
	}
}

// GoldenBraidInputFlags returns the flag for the goldenbraid plasmid CSV input.
// Reuses --input-source, --s3-bucket, --s3-bucket-path from InputFlags() for source routing.
func GoldenBraidInputFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "goldenbraid-input",
			Usage: "GoldenBraid plasmid CSV path (local) or filename in bucket",
			Value: "goldenbraid.csv",
		},
	}
}

// LoggingFlags returns flags for logging configuration
func LoggingFlags() []cli.Flag {
	return []cli.Flag{
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
}
