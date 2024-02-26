package database

import (
	"context"
	"net/http"

	A "github.com/IBM/fp-go/array"
	R "github.com/IBM/fp-go/context/readerioeither"
	H "github.com/IBM/fp-go/context/readerioeither/http"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
	"github.com/dictyBase/modware-import/internal/baserow/client"
	"github.com/dictyBase/modware-import/internal/baserow/httpapi"
)

type OntologyTableManager struct {
	*TableManager
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

func (ont *OntologyTableManager) onFieldDelReqFeedbackSome(
	field tableFieldsResponse,
) fieldsReqFeedback {
	resp := F.Pipe3(
		ont.TableFieldsChangeURL(field),
		F.Bind13of3(H.MakeRequest)("DELETE", nil),
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

func (ont *OntologyTableManager) RemoveField(
	tbl *client.Table, field string,
) (string, error) {
	var empty string
	fields, err := ont.ListTableFields(tbl)
	if err != nil {
		return empty, err
	}
	delOutput := F.Pipe2(
		fields,
		A.FindFirst(HasField(field)),
		O.Fold[tableFieldsResponse](
			onFieldDelReqFeedbackNone,
			ont.onFieldDelReqFeedbackSome,
		),
	)
	if delOutput.Error != nil {
		return empty, delOutput.Error
	}

	return delOutput.Msg, nil
}

func (ont *OntologyTableManager) FieldNames() []string {
	return []string{"term_id", "name", "is_obsolete"}
}

func (ont *OntologyTableManager) FieldDefs() []map[string]interface{} {
	return []map[string]interface{}{
		{"name": "name", "type": "text"},
		{"name": "term_id", "type": "text"},
		{"name": "is_obsolete", "type": "boolean"},
	}
}

func onFieldDelReqFeedbackSuccess(
	resp tableFieldDelResponse,
) fieldsReqFeedback {
	return fieldsReqFeedback{Msg: "deleted extra field"}
}

func onFieldDelReqFeedbackNone() fieldsReqFeedback {
	return fieldsReqFeedback{Msg: "no field found to delete"}
}
