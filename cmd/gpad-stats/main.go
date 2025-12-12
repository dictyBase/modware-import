package main

import (
	"fmt"
	"log"
	"os"

	E "github.com/IBM/fp-go/either"
	"github.com/dictyBase/modware-import/internal/gpadstats"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "gpad-stats",
		Usage: "Generate statistics from Gene Ontology annotation files (GPAD)",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "Path to the GPAD/TSV file",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			filePath := c.String("file")

			// Execute the pure functional pipeline
			result := gpadstats.Run(filePath)

			// Unwrap and handle the result (Impure boundary)
			return E.Fold(
				func(err error) error {
					return cli.Exit(fmt.Sprintf("Error: %v", err), 1)
				},
				func(count int) error {
					fmt.Printf("Unique Gene Count: %d\n", count)
					return nil
				},
			)(result())
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
