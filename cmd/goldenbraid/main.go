package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/dictyBase/modware-import/internal/cli/stockcenter"
	"github.com/dictyBase/modware-import/internal/k8s/wait"
	loader "github.com/dictyBase/modware-import/internal/load/stockcenter"
	"github.com/urfave/cli/v2"
)

func main() { //nolint:funlen
	app := &cli.App{
		Name:  "goldenbraid",
		Usage: "GoldenBraid data loading tools",
		Commands: []*cli.Command{
			{
				Name:   "inventory",
				Usage:  "Load GoldenBraid inventory",
				Action: stockcenter.LoadGoldenBraidInventory,
				Before: stockcenter.SetAllClients,
				Flags: slices.Concat(
					stockcenter.InventoryInputFlags(),
					stockcenter.GoldenBraidInputFlags(),
					stockcenter.S3Flags(),
					stockcenter.StockConnectionFlags(),
					stockcenter.AnnotationFlags(),
					stockcenter.LoggingFlags(),
				),
			},
			{
				Name:   "plasmid",
				Usage:  "Load GoldenBraid plasmid CSV data",
				Action: loader.LoadGoldenBraidCli,
				Before: stockcenter.SetStockAndS3Clients,
				Flags: slices.Concat(
					stockcenter.InputFlags(),
					stockcenter.GoldenBraidInputFlags(),
					stockcenter.S3Flags(),
					stockcenter.StockConnectionFlags(),
					stockcenter.LoggingFlags(),
					[]cli.Flag{
						&cli.StringFlag{
							Name:     "user-email",
							Aliases:  []string{"u"},
							Usage:    "Email of the user loading the data",
							Required: true,
						},
						&cli.StringFlag{
							Name:    "plasmid-cvterm",
							Aliases: []string{"c"},
							Usage:   "Plasmid ontology term",
							Value:   "GB vector",
						},
						&cli.StringFlag{
							Name:  "depositor",
							Usage: "Email of the depositor",
							Value: "gadi@bcm.edu",
						},
					},
				),
			},
			{
				Name:   "plasmid-ontology",
				Usage:  "Update plasmids with ontology term",
				Action: loader.LoadPlasmidOntologyCli,
				Before: stockcenter.SetStockClientWrapper,
				Flags: slices.Concat(
					stockcenter.StockConnectionFlags(),
					stockcenter.OntologyFlags(),
					stockcenter.LoggingFlags(),
				),
			},
			{
				Name:   "wait-job",
				Usage:  "Wait for a Kubernetes job to complete, detecting stuck pods early",
				Action: wait.JobAction,
				Flags: slices.Concat(
					stockcenter.LoggingFlags(),
					[]cli.Flag{
						&cli.StringFlag{
							Name:     "name",
							Usage:    "Job name to wait for",
							Required: true,
						},
						&cli.StringFlag{
							Name:  "namespace",
							Usage: "Kubernetes namespace",
							Value: "dev",
						},
						&cli.StringFlag{
							Name:  "timeout",
							Usage: "Maximum wait duration (e.g. 60s, 5m)",
							Value: "60s",
						},
						&cli.StringFlag{
							Name:    "kubeconfig",
							Usage:   "Path to kubeconfig file",
							EnvVars: []string{"KUBECONFIG"},
						},
					},
				),
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
