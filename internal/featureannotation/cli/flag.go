package cli

import (
	"slices"

	"github.com/urfave/cli/v2"
)

// GeneProductFromCsvFlag returns all flags required for loading gene products from CSV files.
// It combines gene product specific flags with gRPC connection flags.
func GeneProductFromCsvFlag() []cli.Flag {
	return slices.Concat(
		LoadGeneProductFlag(),
		featureAnnotationGrpcFlags(),
	)
}

// GeneDescriptionFromCsvFlag returns all flags required for loading gene descriptions from CSV files.
// It combines gene description specific flags with gRPC connection flags.
func GeneDescriptionFromCsvFlag() []cli.Flag {
	return slices.Concat(
		LoadGeneDescriptionFlag(),
		featureAnnotationGrpcFlags(),
	)
}

// LoadGeneProductFlag returns flags specific to gene product loading operations.
// Input validation is performed during command execution to ensure files exist
// and user email is properly formatted.
func LoadGeneProductFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:     "input",
			Aliases:  []string{"i"},
			Usage:    "one or more input CSV files with gene products (must exist and be readable)",
			Required: true,
		},
		&cli.IntFlag{
			Name:  "workers",
			Usage: "number of concurrent workers for loading (1-50)",
			Value: 4,
		},
		&cli.IntFlag{
			Name:  "batch-size",
			Usage: "batch size for loading (1-1000)",
			Value: 100,
		},
		&cli.StringFlag{
			Name:     "user",
			Usage:    "email address of the user running the load",
			Required: true,
		},
	}
}

// LoadGeneDescriptionFlag returns flags specific to gene description loading operations.
// Input validation is performed during command execution to ensure the file exists
// and user email is properly formatted.
func LoadGeneDescriptionFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "input",
			Aliases:  []string{"i"},
			Usage:    "input CSV file with gene descriptions (must exist and be readable)",
			Required: true,
		},
		&cli.IntFlag{
			Name:  "workers",
			Usage: "number of concurrent workers for loading (1-50)",
			Value: 4,
		},
		&cli.IntFlag{
			Name:  "batch-size",
			Usage: "batch size for loading (1-1000)",
			Value: 100,
		},
		&cli.StringFlag{
			Name:     "user",
			Usage:    "email address of the user running the load",
			Required: true,
		},
	}
}

// arangoDBConnectionFlags returns a common set of ArangoDB connection flags.
// These flags are used across multiple commands that need to connect to ArangoDB.
func arangoDBConnectionFlags() []cli.Flag {
	return []cli.Flag{
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
			Value:   "arangodb",
		},
		&cli.IntFlag{
			Name:    "arangodb-port",
			Usage:   "ArangoDB port",
			EnvVars: []string{"ARANGODB_SERVICE_PORT"},
			Value:   8529,
		},
		&cli.BoolFlag{
			Name:    "is-secure",
			Usage:   "Whether to use TLS for ArangoDB connection",
			EnvVars: []string{"ARANGODB_IS_SECURE"},
			Value:   false,
		},
	}
}

// featureAnnotationGrpcFlags returns gRPC flags for the feature annotation service.
// These flags configure the connection to the feature annotation gRPC API.
func featureAnnotationGrpcFlags() []cli.Flag {
	return []cli.Flag{
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

// LoadFeatureAnnotationFlag returns all flags required for loading feature annotations.
// It combines ArangoDB connection flags, gRPC flags, and pubmed/grpc worker configuration.
func LoadFeatureAnnotationFlag() []cli.Flag {
	return slices.Concat(
		arangoDBConnectionFlags(),
		featureAnnotationGrpcFlags(),
		[]cli.Flag{
			&cli.IntFlag{
				Name:    "pubmed-workers",
				Value:   4,
				Usage:   "Number of pubmed fetcher workers",
				EnvVars: []string{"PUBMED_WORKERS"},
			},
			&cli.IntFlag{
				Name:    "grpc-workers",
				Value:   8,
				Usage:   "Number of gRPC create workers",
				EnvVars: []string{"GRPC_WORKERS"},
			},
		},
	)
}

// LoadCSVToArangodbFlag returns all flags required for loading CSV data to ArangoDB.
// It combines ArangoDB connection flags with CSV processing specific flags.
func LoadCSVToArangodbFlag() []cli.Flag {
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
	return slices.Concat(arangoDBConnectionFlags(), csvFlags)
}

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
	}
	return slices.Concat(
		arangoDBConnectionFlags(),
		featureAnnotationGrpcFlags(),
		geneUpdaterSpecificFlags,
	)
}

// GeneProductUpdaterFlags returns flags for gene product updater
func GeneProductUpdaterFlags() []cli.Flag {
	geneProductSpecificFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "legacy-database",
			Value:   "cgm_ddb",
			Usage:   "Legacy database name to query gene products from",
			EnvVars: []string{"LEGACY_DATABASE"},
		},
		&cli.IntFlag{
			Name:    "legacy-workers",
			Value:   4,
			Usage:   "Number of legacy database query workers",
			EnvVars: []string{"LEGACY_WORKERS"},
		},
		&cli.IntFlag{
			Name:    "grpc-workers", // This flag was already in GeneUpdaterFlags, ensure consistency or rename if needed
			Value:   8,
			Usage:   "Number of gRPC update workers",
			EnvVars: []string{"GRPC_WORKERS"},
		},
	}

	return slices.Concat(
		arangoDBConnectionFlags(),    // Common ArangoDB flags (for the main DB, not legacy)
		featureAnnotationGrpcFlags(), // Common gRPC flags
		geneProductSpecificFlags,     // Specific flags for this updater
	)
}

// SynonymLoaderFlags returns all flags required for the synonym loader command.
func SynonymLoaderFlags() []cli.Flag {
	return slices.Concat(
		arangoDBConnectionFlags(),
		featureAnnotationGrpcFlags(),
		[]cli.Flag{
			&cli.IntFlag{
				Name:    "grpc-workers",
				Value:   4,
				Usage:   "Number of gRPC update workers",
				EnvVars: []string{"GRPC_WORKERS"},
			},
		},
	)
}

// ParseUnknowmeDataFlags returns flags for parsing unknowme HTML data to CSV files.
// Input validation is performed during command execution to ensure files exist.
func ParseUnknowmeDataFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "input",
			Aliases:  []string{"i"},
			Usage:    "input HTML file to parse (must exist and be readable)",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "gene-product-output",
			Aliases: []string{"p"},
			Usage:   "output CSV file for gene products",
			Value:   "gene_products.csv",
		},
		&cli.StringFlag{
			Name:    "gene-description-output",
			Aliases: []string{"d"},
			Usage:   "output CSV file for gene descriptions",
			Value:   "gene_descriptions.csv",
		},
	}
}
