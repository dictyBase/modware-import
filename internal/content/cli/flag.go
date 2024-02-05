package cli

import "github.com/urfave/cli/v2"

func ContentLoaderFlags() []cli.Flag {
	return []cli.Flag{
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
			Name:    "content-grpc-host",
			Usage:   "grpc host address for content service",
			Value:   "stock-api",
			EnvVars: []string{"CONTENT_API_SERVICE_HOST"},
		},
		&cli.StringFlag{
			Name:    "content-grpc-port",
			Usage:   "grpc port for content service",
			EnvVars: []string{"CONTENT_API_SERVICE_PORT"},
		},
	}
}
