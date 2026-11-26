// Command phenoannotation prepares phenotype annotation batches.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v3"

	pheno "github.com/dictyBase/modware-import/internal/phenoannotation/cli"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	app := &cli.Command{
		Name:  "phenoannotation",
		Usage: "Phenotype annotation loading tools",
		Commands: []*cli.Command{
			{
				Name:   "normalize",
				Usage:  "Convert create workbooks into reviewable CSV files",
				Action: pheno.RunNormalize,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "folder",
						Usage:    "Folder of create workbooks",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "out",
						Usage: "New output directory",
						Value: "pheno_normalize",
					},
					&cli.StringFlag{
						Name:  "sheet",
						Usage: "Worksheet name (default: first worksheet)",
					},
				},
			},
		},
	}
	err := app.Run(ctx, os.Args)
	stop()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(pheno.ExitCode(err))
	}
}
