package cli

import (
	A "github.com/IBM/fp-go/array"
	"github.com/urfave/cli/v2"
)

func CreatePhenotypeTableFlag() []cli.Flag {
	flags := append(tableCreationFlags(),
		&cli.StringFlag{
			Name:     "table",
			Usage:    "table to create for loading phenotype annotation",
			Required: true,
		},
	)
	return append(flags, phenoOntologyTableFlags()...)
}

func LoadPhenoFolderToTableFlag() []cli.Flag {
	return append(phenoToTableFlag(),
		&cli.StringFlag{
			Name:     "folder",
			Aliases:  []string{"f"},
			Usage:    "folder with excel spreadsheet files with phenotype annotations",
			Required: true,
		},
	)
}

func LoadPhenoToTableFlag() []cli.Flag {
	return append(phenoToTableFlag(),
		&cli.StringFlag{
			Name:     "input",
			Aliases:  []string{"i"},
			Usage:    "input excel spreadsheet file with phenotype annotations",
			Required: true,
		})
}

func phenoToTableFlag() []cli.Flag {
	return A.ArrayConcatAll[cli.Flag](
		tableCreationFlags(),
		spreadsheetFlag(),
		phenoOntologyTableFlags(),
	)
}

func phenoOntologyTableFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "assay-ontology-table",
			Usage: "table name containing assay ontology",
			Value: "assay_ontology",
		},
		&cli.StringFlag{
			Name:  "env-ontology-table",
			Usage: "table name containing environmental ontology",
			Value: "environment_ontology",
		},
		&cli.StringFlag{
			Name:  "phenotype-ontology-table",
			Usage: "table name containing phenotype ontology",
			Value: "phenotype_ontology",
		},
	}
}
