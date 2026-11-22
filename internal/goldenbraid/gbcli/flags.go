// Package gbcli exposes the GoldenBraid export pipeline as a CLI subcommand.
package gbcli

import "github.com/urfave/cli/v2"

// exportFlags holds the export-specific flags. This file has no fp-go
// import, so the CLI flag literals stay outside the struct literal gate.
func exportFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "input",
			Usage:    "Local path or s3://bucket/key of a .dna file, a .zip archive, or a folder (chosen by extension)",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "output",
			Usage:    "Local directory or s3://bucket/prefix for the fasta, genbank, and svg folders",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "s3-region",
			Usage: "AWS region for the S3 client",
			Value: "us-east-1",
		},
		&cli.BoolFlag{
			Name:  "show-primers",
			Usage: "Draw primer_bind features on the SVG map. GenBank output always keeps them",
		},
	}
}
