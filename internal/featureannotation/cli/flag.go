package cli

import "github.com/urfave/cli/v2"
	"slices"

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
			Value: 40,
		},
		&cli.IntFlag{
			Name:  "workers",
			Usage: "Number of concurrent workers for batch processing",
			Value: 4,
		},
	}

	// Combine the flags
	return append(flags, csvFlags...)
// GeneUpdaterFlags returns all flags required for the gene updater command.
func GeneUpdaterFlags() []cli.Flag {
	geneUpdaterSpecificFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "aql-query",
			Value:   DefaultAQLQuery, // This will come from gene_updater.go
			Usage:   "AQL query to fetch gene data",
			EnvVars: []string{"AQL_QUERY"},
		},
		// Worker and timeout flags
		&cli.IntFlag{
			Name:    "processing-workers",
			Value:   4,
			Usage:   "Number of HTML processing workers",
			EnvVars: []string{"PROCESSING_WORKERS"},
		},
		&cli.IntFlag{
			Name:    "grpc-workers",
			Value:   8,
			Usage:   "Number of gRPC update workers",
			EnvVars: []string{"GRPC_WORKERS"},
		},
		&cli.DurationFlag{
			Name:    "grpc-timeout",
			Usage:   "Timeout for gRPC calls",
			EnvVars: []string{"GRPC_TIMEOUT"},
			Value:   30 * time.Second,
		},
		&cli.DurationFlag{
			Name:    "arango-timeout",
			Usage:   "Timeout for ArangoDB query execution",
			EnvVars: []string{"ARANGO_TIMEOUT"},
			Value:   5 * time.Minute,
		},
	}
	return slices.Concat(
		arangoDBConnectionFlags(),
		featureAnnotationGrpcFlags(),
		geneUpdaterSpecificFlags,
	)
}
