package cli

import (
	"github.com/urfave/cli/v2"
)

func CreateStrainTableFlag() []cli.Flag {
	tblFlags := append(tableCreationFlags(),
		&cli.StringFlag{
			Name:     "table",
			Usage:    "table to create for loading strain annotation",
			Required: true,
		},
	)
	return append(tblFlags, strainOntologyTableFlags()...)
}

func LoadStrainFolderToTableFlag() []cli.Flag {
	return append(strainToTableFlag(),
		&cli.StringFlag{
			Name:     "folder",
			Aliases:  []string{"f"},
			Usage:    "folder with excel spreadsheet files with strain annotations",
			Required: true,
		})
}

func LoadStrainToTableFlag() []cli.Flag {
	return append(strainToTableFlag(),
		&cli.StringFlag{
			Name:     "input",
			Aliases:  []string{"i"},
			Usage:    "input excel spreadsheet file with strain annotations",
			Required: true,
		})
}

func strainToTableFlag() []cli.Flag {
	tblFlags := append(tableCreationFlags(), spreadsheetFlag()...)
	return append(tblFlags, strainOntologyTableFlags()...)
}

func strainOntologyTableFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "strainchar-ontology-table",
			Usage: "table name containing strain characteristics ontology",
			Value: "strain_characteristics_ontology",
		},
		&cli.StringFlag{
			Name:  "genetic-mod-ontology-table",
			Usage: "table name containing genetic modification ontology",
			Value: "genetic_modification_ontology",
		},
		&cli.StringFlag{
			Name:  "mutagenesis-method-ontology-table",
			Usage: "table name containing mutagenesis method ontology",
			Value: "mutagenesis_method_ontology",
		},
	}
}
