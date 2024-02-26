package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	A "github.com/IBM/fp-go/array"
	O "github.com/IBM/fp-go/option"

	"github.com/dictyBase/modware-import/internal/collection"
	"golang.org/x/exp/slices"

	"github.com/dictyBase/modware-import/internal/baserow/client"
	"github.com/dictyBase/modware-import/internal/baserow/httpapi"
	"github.com/sirupsen/logrus"

	R "github.com/IBM/fp-go/context/readerioeither"
	H "github.com/IBM/fp-go/context/readerioeither/http"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
)

type tableFieldsResponse struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type TableManager struct {
	Logger     *logrus.Entry
	Client     *client.APIClient
	Ctx        context.Context
	Token      string
	DatabaseId int32
}

var (
	HasField = F.Curry2(uncurriedHasField)
)

func (tbm *TableManager) TableFieldsChangeURL(
	field tableFieldsResponse,
) string {
	return fmt.Sprintf(
		"https://%s/api/database/fields/%d/",
		tbm.Client.GetConfig().Host,
		field.Id,
	)
}

func (tbm *TableManager) TableFieldsURL(tbl *client.Table) string {
	return fmt.Sprintf(
		"https://%s/api/database/fields/table/%d/",
		tbm.Client.GetConfig().Host,
		tbl.GetId(),
	)
}

func (tbm *TableManager) CreateTable(
	table string, fields []string,
) (*client.Table, error) {
	isTrue := true
	var row []interface{}
	row = append(row, fields)
	tbl, resp, err := tbm.Client.
		DatabaseTablesApi.
		CreateDatabaseTable(tbm.Ctx, tbm.DatabaseId).
		TableCreate(client.TableCreate{Name: table, Data: row, FirstRowHeader: &isTrue}).
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

func (tbm *TableManager) TableFieldsResp(
	tbl *client.Table,
) (*http.Response, error) {
	reqURL := fmt.Sprintf(
		"https://%s/api/database/fields/table/%d/",
		tbm.Client.GetConfig().Host,
		tbl.GetId(),
	)
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error in creating request %s ", err)
	}
	httpapi.CommonHeader(req, tbm.Token, "Token")
	return httpapi.ReqToResponse(req)
}

func (tbm *OntologyTableManager) ListTableFields(
	tbl *client.Table,
) ([]tableFieldsResponse, error) {
	resp := F.Pipe3(
		tbm.TableFieldsURL(tbl),
		H.MakeGetRequest,
		R.Map(httpapi.SetHeaderWithJWT(tbm.Token)),
		readFieldsResp,
	)(context.Background())
	output := F.Pipe1(
		resp(),
		E.Fold[error, []tableFieldsResponse, fieldsReqFeedback](
			onFieldsReqFeedbackError,
			onFieldsReqFeedbackSuccess,
		),
	)
	return output.Fields, output.Error
}

func (tbm *OntologyTableManager) CreateFields(tbl *client.Table) error {
	reqURL := fmt.Sprintf(
		"https://%s/api/database/fields/table/%d/",
		tbm.Client.GetConfig().Host,
		tbl.GetId(),
	)
	for _, payload := range tbm.FieldDefs() {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("error in encoding body %s", err)
		}
		req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("error in creating request %s ", err)
		}
		httpapi.CommonHeader(req, tbm.Token, "JWT")
		res, err := httpapi.ReqToResponse(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()
	}

	return nil
}

func (tbm *OntologyTableManager) CheckAllTableFields(
	tbl *client.Table,
) (bool, error) {
	ok := false
	res, err := tbm.TableFieldsResp(tbl)
	if err != nil {
		return ok, err
	}
	defer res.Body.Close()
	existing := make([]tableFieldsResponse, 0)
	if err := json.NewDecoder(res.Body).Decode(&existing); err != nil {
		return ok, fmt.Errorf("error in decoding response %s", err)
	}
	exFields := collection.Map(
		existing,
		func(input tableFieldsResponse) string { return input.Name },
	)
	for _, fld := range tbm.FieldNames() {
		if num := slices.Index(exFields, fld); num == -1 {
			return ok, nil
		}
	}

	return true, nil
}

func (tbm *OntologyTableManager) RemoveField(
	tbl *client.Table, field string,
) (string, error) {
	var empty string
	fields, err := tbm.ListTableFields(tbl)
	if err != nil {
		return empty, err
	}
	delOutput := F.Pipe2(
		fields,
		A.FindFirst(HasField(field)),
		O.Fold[tableFieldsResponse](
			onFieldDelReqFeedbackNone,
			tbm.onFieldDelReqFeedbackSome,
		),
	)
	if delOutput.Error != nil {
		return empty, delOutput.Error
	}

	return delOutput.Msg, nil
}

func uncurriedHasField(name string, fieldResp tableFieldsResponse) bool {
	return fieldResp.Name == name
}

func onFieldsReqFeedbackError(err error) fieldsReqFeedback {
	return fieldsReqFeedback{Error: err}
}

func onFieldsReqFeedbackSuccess(resp []tableFieldsResponse) fieldsReqFeedback {
	return fieldsReqFeedback{Fields: resp}
}
