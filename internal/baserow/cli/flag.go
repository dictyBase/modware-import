package cli

import "github.com/urfave/cli/v2"

func CreateAccessTokenFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "email",
			Aliases:  []string{"e"},
			Usage:    "Email of the user",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "password",
			Aliases:  []string{"p"},
			Usage:    "Database password",
			Required: true,
		},
		&cli.BoolFlag{
			Name:  "save-refresh-token",
			Usage: "whether to persist the refresh token",
			Value: true,
		},
		&cli.StringFlag{
			Name:  "refresh-token-path",
			Usage: "where the refresh token will be saved",
			Value: "./refresh-token.txt",
		},
	}
}

func CreateDatabaseTokenFlag() []cli.Flag {
	aflags := CreateAccessTokenFlag()
	return append(aflags, []cli.Flag{
		&cli.StringFlag{
			Name:     "workspace",
			Aliases:  []string{"w"},
			Usage:    "Only tables under this workspaces can be accessed",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "token name",
			Required: true,
		},
	}...)
}

func LoadOntologyToTableFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "token",
			Aliases:  []string{"t"},
			Usage:    "database token with write privilege",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "table-id",
			Usage:    "Database table id",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "input",
			Aliases:  []string{"i"},
			Usage:    "input json formatted ontology file",
			Required: true,
		},
	}
}

func CreateStrainTableFlag() []cli.Flag {
	return append(tableCreationFlags(),
		&cli.StringFlag{
			Name:     "table",
			Usage:    "table to create for loading phenotype annotation",
			Required: true,
		},
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
	)
}

func CreatePhenotypeTableFlag() []cli.Flag {
	return append(tableCreationFlags(),
		&cli.StringFlag{
			Name:     "table",
			Usage:    "table to create for loading phenotype annotation",
			Required: true,
		},
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
	)
}

func CreateOntologyTableFlag() []cli.Flag {
	return append(tableCreationFlags(),
		&cli.StringSliceFlag{
			Name:     "table",
			Usage:    "tables to create for loading ontology",
			Required: true,
		},
	)
}

func tableCreationFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "token",
			Aliases: []string{"t"},
			Usage:   "database token with write privilege",
		},
		&cli.StringFlag{
			Name:  "refresh-token-path",
			Usage: "location, in absence of token value the refresh token will be read",
			Value: "./refresh-token.txt",
		},
		&cli.IntFlag{
			Name:     "database-id",
			Usage:    "Database id",
			Required: true,
		},
	}
}
