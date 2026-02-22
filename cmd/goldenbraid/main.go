package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/dictyBase/modware-import/internal/cli/stockcenter"
	loader "github.com/dictyBase/modware-import/internal/load/stockcenter"
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
				Flags: slices.Concat(
					stockcenter.InputFlags(),
					stockcenter.S3Flags(),
					stockcenter.StockConnectionFlags(),
					stockcenter.AnnotationFlags(),
					stockcenter.LoggingFlags(),
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
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
