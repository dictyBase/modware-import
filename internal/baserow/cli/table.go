package cli

import (
	"context"
	"fmt"

	"github.com/dictyBase/modware-import/internal/baserow/client"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/urfave/cli/v2"
)

type fieldFn func(string) client.FieldCreateField

type CreateFieldProperties struct {
	Client    *client.APIClient
	Ctx       context.Context
	TableId   int
	Field     string
	FieldType client.Type712Enum
}

func MapFieldTypeToFn() map[client.Type712Enum]fieldFn {
	fieldFnMap := make(map[client.Type712Enum]fieldFn)
	fieldFnMap[client.BOOLEAN] = func(field string) client.FieldCreateField {
		return client.FieldCreateField{
			BooleanFieldCreateField: client.NewBooleanFieldCreateField(
				field,
				client.BOOLEAN,
			)}
	}
	fieldFnMap[client.TEXT] = func(field string) client.FieldCreateField {
		return client.FieldCreateField{
			TextFieldCreateField: client.NewTextFieldCreateField(
				field,
				client.TEXT,
			)}
	}
	return fieldFnMap
}

func CreateTableField(args *CreateFieldProperties) error {
	mapper := MapFieldTypeToFn()
	if _, ok := mapper[args.FieldType]; !ok {
		return fmt.Errorf("cannot find field type %s", args.FieldType)
	}
	createFn := mapper[args.FieldType]
	_, resp, err := args.Client.
		DatabaseTableFieldsApi.
		CreateDatabaseTableField(args.Ctx, int32(args.TableId)).
		FieldCreateField(createFn(args.Field)).
		Execute()
	if err != nil {
		return fmt.Errorf("error in creating field %s %s", args.Field, err)
	}
	defer resp.Body.Close()
	return nil
}

func LoadOntologyToTable(cltx *cli.Context) error {
	logger := registry.GetLogger()
	bclient := baserowClient(cltx.String("server"))
	authCtx := context.WithValue(
		context.Background(),
		client.ContextDatabaseToken,
		cltx.String("token"),
	)
	fieldMap := map[string]client.Type712Enum{
		"Name":        client.TEXT,
		"Id":          client.TEXT,
		"Is_obsolete": client.BOOLEAN,
	}
	tlist, resp, err := bclient.
		DatabaseTableFieldsApi.
		ListDatabaseTableFields(authCtx, int32(cltx.Int("table-id"))).
		Execute()
	if err != nil {
		return cli.Exit(
			fmt.Sprintf("error in getting list of table fields %s", err),
			2,
		)
	}
	defer resp.Body.Close()
	if len(tlist) != 0 {
		logger.Debug("fields exists in the database table")
		return nil
	}
	logger.Debug("need to create fields in the table")
	for field, fieldType := range fieldMap {
		err := CreateTableField(&CreateFieldProperties{
			Client:    bclient,
			Ctx:       authCtx,
			TableId:   cltx.Int("table-id"),
			Field:     field,
			FieldType: fieldType,
		})
		if err != nil {
			return cli.Exit(err.Error(), 2)
		}
		logger.Infof("created field %s", field)
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
	tbl, resp, err := bclient.
		DatabaseTablesApi.
		CreateDatabaseTable(authCtx, int32(c.Int("database-id"))).
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
