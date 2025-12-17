package main

import (
	"log"
	"os"

	"github.com/dictyBase/modware-import/internal/cli/stockcenter"
	"github.com/urfave/cli/v2"
)

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
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "input",
						Usage:    "Input CSV file path",
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
						Value: "import/stockcenter",
					},
					&cli.StringFlag{
						Name:  "stock-grpc-host",
						Usage: "gRPC host address for stock service",
						Value: "stock-api",
					},
					&cli.StringFlag{
						Name:  "stock-grpc-port",
						Usage: "gRPC port for stock service",
						Value: "9560",
					},
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
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
