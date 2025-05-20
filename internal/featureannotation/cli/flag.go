package cli

import "github.com/urfave/cli/v2"

// LoadFeatureAnnotationFlag returns all flags required for loading feature annotations
func LoadFeatureAnnotationFlag() []cli.Flag {
	return []cli.Flag{
		// ArangoDB connection flags
		&cli.StringFlag{
			Name:     "arangodb-user",
			Usage:    "ArangoDB user name",
			EnvVars:  []string{"ARANGODB_USER"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "arangodb-pass",
			Usage:    "ArangoDB password",
			EnvVars:  []string{"ARANGODB_PASS"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "arangodb-database",
			Usage:    "ArangoDB database name",
			EnvVars:  []string{"ARANGODB_DATABASE"},
			Required: true,
		},
		&cli.StringFlag{
			Name:    "arangodb-host",
			Usage:   "ArangoDB host",
			EnvVars: []string{"ARANGODB_SERVICE_HOST"},
		},
		&cli.IntFlag{
			Name:    "arangodb-port",
			Usage:   "ArangoDB port",
			EnvVars: []string{"ARANGODB_SERVICE_PORT"},
		},
		&cli.BoolFlag{
			Name:    "is-secure",
			Usage:   "Whether to use TLS for ArangoDB connection",
			EnvVars: []string{"ARANGODB_IS_SECURE"},
		},

		// Feature annotation gRPC service flags
		&cli.StringFlag{
			Name:    "feature-annotation-grpc-host",
			Usage:   "Feature annotation gRPC host",
			EnvVars: []string{"ANNO_FEAT_API_SERVICE_HOST"},
			Value:   "anno-feat-api",
		},
		&cli.StringFlag{
			Name:    "feature-annotation-grpc-port",
			Usage:   "Feature annotation gRPC port",
			EnvVars: []string{"ANNO_FEAT_API_SERVICE_PORT"},
			Value:   "9250",
		},
	}
}

// LoadCSVToArangodbFlag returns all flags required for loading CSV data to ArangoDB
func LoadCSVToArangodbFlag() []cli.Flag {
	// Reuse existing ArangoDB connection flags
	// Assuming the first 6 flags are ArangoDB related based on the
	// LoadFeatureAnnotationFlag definition
	flags := LoadFeatureAnnotationFlag()[:6]

	csvFlags := []cli.Flag{
		&cli.StringFlag{
			Name:     "csv-file",
			Usage:    "Path to CSV file to load",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "collection",
			Usage: "ArangoDB collection name to update",
			Value: "featureprop",
		},
		&cli.StringFlag{
			Name:  "delimiter",
			Usage: "CSV delimiter character",
			Value: ",",
		},
		&cli.IntFlag{
			Name:  "batch-size",
			Usage: "Number of documents to update in a single batch",
			Value: 500,
		},
		&cli.IntFlag{
			Name:  "workers",
			Usage: "Number of concurrent workers for batch processing",
			Value: 4,
		},
		&cli.BoolFlag{
			Name:  "concurrent",
			Usage: "Whether to use the concurrent worker pool implementation",
			Value: false,
		},
	}

	// Combine the flags
	return append(flags, csvFlags...)
}
