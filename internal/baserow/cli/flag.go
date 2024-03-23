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

func spreadsheetFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "sheet",
			Aliases: []string{"s"},
			Usage:   "name of sheet which contains the annotation",
			Value:   "Strain_Annotations",
		},
		&cli.IntFlag{
			Name:     "table-id",
			Usage:    "Database table id",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "workspace",
			Aliases:  []string{"w"},
			Usage:    "name of the workspace whether the database exists",
			Required: true,
		},
	}
}
