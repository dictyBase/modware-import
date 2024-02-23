package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/dictyBase/modware-import/internal/baserow/client"
	"github.com/dictyBase/modware-import/internal/baserow/httpapi"
	"github.com/sirupsen/logrus"
)

type OntologyTableManager struct {
	Logger     *logrus.Entry
	Client     *client.APIClient
	Ctx        context.Context
	Token      string
	DatabaseId int32
}

func (ont *OntologyTableManager) CreateTable(
	table string,
) (*client.Table, error) {
	var row []interface{}
	row = append(row, []string{""})
	tbl, resp, err := ont.Client.
		DatabaseTablesApi.
		CreateDatabaseTable(ont.Ctx, ont.DatabaseId).
		TableCreate(client.TableCreate{Name: table, Data: row}).
		Execute()
	if err != nil {
		return tbl, fmt.Errorf(
			"error in creating table %s %s",
			table, err,
		)
	}
	defer resp.Body.Close()

	return tbl, nil
}

func (ont *OntologyTableManager) CheckAllTableFields(
	tableId int32,
	fields []string,
) (bool, error) {
	ok := false
	tlist, resp, err := ont.Client.
		DatabaseTableFieldsApi.
		ListDatabaseTableFields(ont.Ctx, tableId).
		Execute()
	if err != nil {
		return ok, fmt.Errorf(
			"error in getting list of table fields %s",
			err,
		)
	}
	defer resp.Body.Close()
	existFields := toFieldNames(tlist)
	slices.Sort(existFields)
	slices.Sort(fields)
	return slices.Equal(
		fields,
		toFieldNames(tlist),
	), nil
}

func (ont *OntologyTableManager) RemoveInitialFields(tbl *client.Table) error {
	return nil
}

func (ont *OntologyTableManager) CreateFields(tbl *client.Table) error {
	reqURL := fmt.Sprintf(
		"https://%s/api/database/fields/table/%d/",
		ont.Client.GetConfig().Host,
		tbl.GetId(),
	)
	for _, payload := range fieldDefs() {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("error in encoding body %s", err)
		}
		req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("error in creating request %s ", err)
		}
		httpapi.CommonHeader(req, ont.Token, "JWT")
		res, err := httpapi.ReqToResponse(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()
	}

	return nil
}

func fieldInformation(field client.FieldField) (int32, string) {
	if field.TextFieldField != nil {
		return field.TextFieldField.Id, field.TextFieldField.Name
	}
	if field.BooleanFieldField != nil {
		return field.BooleanFieldField.Id, field.BooleanFieldField.Name
	}
	if field.LongTextFieldField != nil {
		return field.LongTextFieldField.Id, field.LongTextFieldField.Name
	}
	return field.DateFieldField.Id, *field.DateFieldField.Name
}

func toFieldNames(fields []client.FieldField) []string {
	fieldNames := make([]string, 0)
	for _, fld := range fields {
		_, name := fieldInformation(fld)
		fieldNames = append(fieldNames, name)
	}
	return fieldNames
}

func fieldDefs() []map[string]interface{} {
	return []map[string]interface{}{
		{"name": "name", "type": "text"},
		{"name": "term_id", "type": "text"},
		{"name": "is_obsolete", "type": "boolean"},
	}
}
