package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	A "github.com/IBM/fp-go/array"
	J "github.com/IBM/fp-go/json"
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

type FieldDefinition interface {
	FieldNames() []string
	FieldDefs() map[string]interface{}
}

type TableManager struct {
	FieldDefinition
	Logger     *logrus.Entry
	Client     *client.APIClient
	Ctx        context.Context
	Token      string
	DatabaseId int32
}

func (tbm *TableManager) TableFieldsChangeURL(
	req tableFieldReq,
) string {
	return fmt.Sprintf(
		"https://%s/api/database/fields/%d/",
		tbm.Client.GetConfig().Host,
		req.Id,
	)
}

func (tbm *TableManager) TableFieldsURL(tbl *client.Table) string {
	return fmt.Sprintf(
		"https://%s/api/database/fields/table/%d/",
		tbm.Client.GetConfig().Host,
		tbl.GetId(),
	)
}

func (tbm *TableManager) CreateTableURL() string {
	return fmt.Sprintf(
		"https://%s/api/database/tables/database/%d/",
		tbm.Client.GetConfig().Host,
		tbm.DatabaseId,
	)
}

func (tbm *TableManager) CreateTable(
	table string, fields []string,
) (*client.Table, error) {
	var row []interface{}
	params := map[string]interface{}{
		"name":             table,
		"data":             append(row, fields),
		"first_row_header": "true",
	}
	createPayload := F.Pipe2(
		params,
		J.Marshal,
		E.Fold(onJSONPayloadError, onJSONPayloadSuccess),
	)
	if createPayload.Error != nil {
		return &client.Table{}, createPayload.Error
	}
	resp := F.Pipe3(
		tbm.CreateTableURL(),
		makeHTTPRequest("POST", bytes.NewBuffer(createPayload.Payload)),
		R.Map(httpapi.SetHeaderWithJWT(tbm.Token)),
		readTableCreateResp,
	)(context.Background())
	output := F.Pipe1(
		resp(),
		E.Fold[error, tableFieldRes, fieldsReqFeedback](
			onFieldsReqFeedbackError,
			onTableCreateFeedbackSuccess,
		),
	)

	return output.Table, output.Error
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

func (tbm *TableManager) ListTableFields(
	tbl *client.Table,
) ([]tableFieldRes, error) {
	resp := F.Pipe3(
		tbm.TableFieldsURL(tbl),
		H.MakeGetRequest,
		R.Map(httpapi.SetHeaderWithJWT(tbm.Token)),
		readFieldsResp,
	)(context.Background())
	output := F.Pipe1(
		resp(),
		E.Fold[error, []tableFieldRes, fieldsReqFeedback](
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
	for _, params := range tbm.FieldDefs() {
		jsonData, err := json.Marshal(params)
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
	existing := make([]tableFieldRes, 0)
	if err := json.NewDecoder(res.Body).Decode(&existing); err != nil {
		return ok, fmt.Errorf("error in decoding response %s", err)
	}
	exFields := collection.Map(
		existing,
		func(input tableFieldRes) string { return input.Name },
	)
	for _, fld := range tbm.FieldNames() {
		if num := slices.Index(exFields, fld); num == -1 {
			return ok, nil
		}
	}

	return true, nil
}

func (tbm *TableManager) UpdateField(
	tbl *client.Table,
	field string,
	updateSpec map[string]interface{},
) (string, error) {
	var empty string
	fields, err := tbm.ListTableFields(tbl)
	if err != nil {
		return empty, err
	}
	updateOutput := F.Pipe3(
		fields,
		A.FindFirst(HasField(field)),
		O.Map(ResToReqTableWithParams(updateSpec)),
		O.Fold[tableFieldReq](
			onFieldDelReqFeedbackNone,
			tbm.onFieldUpdateReqFeedbackSome,
		),
	)

	return updateOutput.Msg, updateOutput.Error
}

func (tbm *TableManager) RemoveField(
	tbl *client.Table, req string,
) (string, error) {
	var empty string
	fields, err := tbm.ListTableFields(tbl)
	if err != nil {
		return empty, err
	}
	delOutput := F.Pipe3(
		fields,
		A.FindFirst(HasField(req)),
		O.Map(ResToReqTable),
		O.Fold[tableFieldReq](
			onFieldDelReqFeedbackNone,
			tbm.onFieldDelReqFeedbackSome,
		),
	)

	return delOutput.Msg, delOutput.Error
}

func (tbm *TableManager) onFieldUpdateReqFeedbackSome(
	req tableFieldReq,
) fieldsReqFeedback {
	payloadResp := F.Pipe2(
		req.Params,
		J.Marshal,
		E.Fold(onJSONPayloadError, onJSONPayloadSuccess),
	)
	if payloadResp.Error != nil {
		return fieldsReqFeedback{Error: payloadResp.Error}
	}
	resp := F.Pipe3(
		tbm.TableFieldsChangeURL(req),
		makeHTTPRequest("PATCH", bytes.NewBuffer(payloadResp.Payload)),
		R.Map(httpapi.SetHeaderWithJWT(tbm.Token)),
		readUpdateFieldsResp,
	)(context.Background())

	return F.Pipe1(
		resp(),
		E.Fold[error, tableFieldUpdateResponse, fieldsReqFeedback](
			onFieldsReqFeedbackError,
			onFieldUpdateReqFeedbackSuccess,
		),
	)
}

func (ont *TableManager) onFieldDelReqFeedbackSome(
	req tableFieldReq,
) fieldsReqFeedback {
	resp := F.Pipe3(
		ont.TableFieldsChangeURL(req),
		makeHTTPRequest("DELETE", nil),
		R.Map(httpapi.SetHeaderWithJWT(ont.Token)),
		readFieldDelResp,
	)(context.Background())

	return F.Pipe1(
		resp(),
		E.Fold[error, tableFieldDelResponse, fieldsReqFeedback](
			onFieldsReqFeedbackError,
			onFieldDelReqFeedbackSuccess,
		),
	)
}
