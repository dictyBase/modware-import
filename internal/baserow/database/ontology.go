package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"slices"

	A "github.com/IBM/fp-go/array"
	R "github.com/IBM/fp-go/context/readerioeither"
	H "github.com/IBM/fp-go/context/readerioeither/http"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
	"github.com/dictyBase/modware-import/internal/baserow/client"
	"github.com/dictyBase/modware-import/internal/baserow/httpapi"
	"github.com/dictyBase/modware-import/internal/collection"
)

type OntologyTableManager struct {
	*TableManager
}

type tableFieldsResponse struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type tableFieldDelResponse struct {
	RelatedFields []struct {
		ID      int `json:"id"`
		TableID int `json:"table_id"`
	} `json:"related_fields,omitempty"`
}

type fieldsReqFeedback struct {
	Error  error
	Fields []tableFieldsResponse
	Msg    string
}

var (
	readFieldDelResp = H.ReadJson[tableFieldDelResponse](
		H.MakeClient(http.DefaultClient),
	)
	readFieldsResp = H.ReadJson[[]tableFieldsResponse](
		H.MakeClient(http.DefaultClient),
	)
)

func (ont *OntologyTableManager) CreateFields(tbl *client.Table) error {
	reqURL := fmt.Sprintf(
		"https://%s/api/database/fields/table/%d/",
		ont.Client.GetConfig().Host,
		tbl.GetId(),
	)
	for _, payload := range FieldDefs() {
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

func (ont *OntologyTableManager) CheckAllTableFields(
	tbl *client.Table,
) (bool, error) {
	ok := false
	res, err := ont.TableFieldsResp(tbl)
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
	for _, fld := range FieldNames() {
		if num := slices.Index(exFields, fld); num == -1 {
			return ok, nil
		}
	}

	return true, nil
}

func FieldNames() []string {
	return []string{"term_id", "name", "is_obsolete"}
}

func FieldDefs() []map[string]interface{} {
	return []map[string]interface{}{
		{"name": "name", "type": "text"},
		{"name": "term_id", "type": "text"},
		{"name": "is_obsolete", "type": "boolean"},
	}
}

func hasExtraField(elem tableFieldsResponse) bool {
	return elem.Name == "Field 1"
}

func onFieldsReqFeedbackError(err error) fieldsReqFeedback {
	return fieldsReqFeedback{Error: err}
}

func onFieldDelReqFeedbackSuccess(
	resp tableFieldDelResponse,
) fieldsReqFeedback {
	return fieldsReqFeedback{Msg: "deleted extra field"}
}

func onFieldsReqFeedbackSuccess(resp []tableFieldsResponse) fieldsReqFeedback {
	return fieldsReqFeedback{Fields: resp}
}

func onFieldDelReqFeedbackNone() fieldsReqFeedback {
	return fieldsReqFeedback{Msg: "no field found to delete"}
}
