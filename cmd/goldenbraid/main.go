package main

import (
	"fmt"
	"os"

	A "github.com/IBM/fp-go/v2/array"
	"github.com/dictyBase/modware-import/internal/goldenbraid/gbcli"

	"github.com/dictyBase/modware-import/internal/cli/stockcenter"
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
				Flags: A.ArrayConcatAll(
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
				Flags: A.ArrayConcatAll(
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
				Flags: A.ArrayConcatAll(
					stockcenter.StockConnectionFlags(),
					stockcenter.OntologyFlags(),
					stockcenter.LoggingFlags(),
				),
			},
			{
				Name:   "export",
				Usage:  "Export GoldenBraid assemblies to FASTA, GenBank, and SVG",
				Action: gbcli.RunExport,
				Flags:  gbcli.ExportFlags(),
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
