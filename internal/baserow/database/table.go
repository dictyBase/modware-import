package database

import (
	"context"
	"fmt"

	"github.com/dictyBase/modware-import/internal/baserow/client"
	"github.com/sirupsen/logrus"
)

type fieldFn func(string) client.FieldCreateField

type CreateFieldProperties struct {
	Client    *client.APIClient
	Ctx       context.Context
	TableId   int
	Field     string
	FieldType client.Type712Enum
}

type OntologyTableFieldsProperties struct {
	Client   *client.APIClient
	Ctx      context.Context
	Logger   *logrus.Entry
	TableId  int
	FieldMap map[string]client.Type712Enum
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

func CreateOntologyTableFields(
	args *OntologyTableFieldsProperties,
) (bool, error) {
	logger := args.Logger
	bclient := args.Client
	authCtx := args.Ctx
	ok := false
	tlist, resp, err := bclient.
		DatabaseTableFieldsApi.
		ListDatabaseTableFields(authCtx, int32(args.TableId)).
		Execute()
	if err != nil {
		return ok, fmt.Errorf("error in getting list of table fields %s", err)
	}
	defer resp.Body.Close()
	if len(tlist) != 0 {
		logger.Debug("fields exists in the database table")
		return ok, nil
	}
	logger.Debug("need to create fields in the table")
	ok = true
	for field, fieldType := range args.FieldMap {
		err := CreateTableField(&CreateFieldProperties{
			Client:    bclient,
			Ctx:       authCtx,
			TableId:   args.TableId,
			Field:     field,
			FieldType: fieldType,
		})
		if err != nil {
			return ok, err
		}
		logger.Infof("created field %s", field)
	}

	return ok, nil
}
