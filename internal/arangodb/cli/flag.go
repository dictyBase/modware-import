package cli

import "github.com/urfave/cli/v2"

func ArangodbLoaderFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "arangodb-host",
			Usage:   "ArangoDB server host",
			EnvVars: []string{"ARANGODB_SERVICE_HOST"},
		},
		&cli.StringFlag{
			Name:  "arangodb-port",
			Usage: "ArangoDB server port",
			Value: "8529",
		},
		&cli.StringFlag{
			Name:     "arangodb-database",
			Usage:    "ArangoDB database name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "arangodb-user",
			Usage:    "ArangoDB user name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "arangodb-pass",
			Usage:    "ArangoDB password",
			Required: true,
		},
		// Reuse existing S3/Minio flags
		&cli.StringFlag{
			Name:    "s3-server-port",
			Usage:   "Port of S3/minio server",
			EnvVars: []string{"MINIO_SERVICE_PORT"},
		},
		&cli.StringFlag{
			Name:    "access-key",
			Usage:   "access key for S3/minio server",
			EnvVars: []string{"ACCESS_KEY"},
		},
		&cli.StringFlag{
			Name:    "secret-key",
			Usage:   "secret key for S3/minio server",
			EnvVars: []string{"SECRET_KEY"},
		},
		&cli.StringFlag{
			Name:    "s3-server",
			Usage:   "S3/minio server endpoint",
			Value:   "minio",
			EnvVars: []string{"MINIO_SERVICE_HOST"},
		},
		&cli.StringFlag{
			Name:  "s3-bucket",
			Usage: "S3/minio bucket for data folder",
			Value: "dictybase",
		},
		&cli.StringFlag{
			Name:     "s3-bucket-path",
			Usage:    "path inside S3 bucket for input files",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "output-dir",
			Usage:    "Directory where JSON files will be written before import",
			Required: true,
		},
	}
}
