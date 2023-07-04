package cli

import (
	"context"
	"fmt"

	"github.com/dictyBase/modware-import/internal/baserow/client"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/urfave/cli/v2"
)

func LoadOntologyToTable(c *cli.Context) error {
	logger := registry.GetLogger()
	bclient := baserowClient(c.String("server"))
	authCtx := context.WithValue(
		context.Background(),
		client.ContextAccessToken,
		c.String("token"),
	)
	tlist, resp, err := bclient.
		DatabaseTableFieldsApi.
		ListDatabaseTableFields(authCtx, int32(c.Int("table-id"))).
		Execute()
	if err != nil {
		return cli.Exit(
			fmt.Errorf("error in getting list of table fields %s", err), 2,
		)
	}
	defer resp.Body.Close()
	if len(tlist) == 0 {
		logger.Debug("need to create fields in the table")
		_, trsp, err := bclient.
			DatabaseTableFieldsApi.
			CreateDatabaseTableField(authCtx, int32(c.Int("table-id"))).
			FieldCreateField(client.FieldCreateField{
				BooleanFieldCreateField: client.NewBooleanFieldCreateField(
					"Is_obsolete",
					client.BOOLEAN,
				),
			}).Execute()
		if err != nil {
			return cli.Exit(
				fmt.Errorf(
					"error in creating table field Is_obsolete %s",
					err,
				),
				2,
			)
		}
		logger.Info("created field Is_obsolete")

		defer trsp.Body.Close()
		for _, field := range []string{"Name", "Id"} {
			_, frsp, err := bclient.
				DatabaseTableFieldsApi.
				CreateDatabaseTableField(authCtx, int32(c.Int("table-id"))).
				FieldCreateField(client.FieldCreateField{
					TextFieldCreateField: client.NewTextFieldCreateField(
						field,
						client.TEXT,
					),
				}).Execute()
			if err != nil {
				return cli.Exit(
					fmt.Errorf(
						"error in creating table field %s %s",
						field, err,
					),
					2,
				)
			}
			defer frsp.Body.Close()
			logger.Infof("created field %s", field)
		}
	}
	return nil
}

func CreateTable(c *cli.Context) error {
	logger := registry.GetLogger()
	bclient := baserowClient(c.String("server"))
	authCtx := context.WithValue(
		context.Background(),
		client.ContextAccessToken,
		c.String("token"),
	)
	tbl, resp, err := bclient.DatabaseTablesApi.CreateDatabaseTable(authCtx, int32(c.Int("database-id"))).
		TableCreate(client.TableCreate{Name: c.String("table")}).
		Execute()
	if err != nil {
		return cli.Exit(
			fmt.Errorf("error in creating table %s", err), 2,
		)
	}
	defer resp.Body.Close()
	logger.Infof("created table %s", tbl.GetName())
	return nil
}

func CreateTableFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "token",
			Aliases:  []string{"t"},
			Usage:    "database token with write privilege",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "database-id",
			Usage:    "Database id",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "table",
			Usage:    "Database table",
			Required: true,
		},
	}
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
	}
}
