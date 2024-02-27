package database

import (
	"net/http"

	H "github.com/IBM/fp-go/context/readerioeither/http"
	F "github.com/IBM/fp-go/function"
	"github.com/dictyBase/modware-import/internal/baserow/client"
)

var (
	makeHTTPRequest  = F.Bind13of3(H.MakeRequest)
	readFieldDelResp = H.ReadJson[tableFieldDelResponse](
		H.MakeClient(http.DefaultClient),
	)
	readFieldsResp = H.ReadJson[[]tableFieldRes](
		H.MakeClient(http.DefaultClient),
	)
	readUpdateFieldsResp = H.ReadJson[tableFieldUpdateResponse](
		H.MakeClient(http.DefaultClient),
	)
	readTableCreateResp = H.ReadJson[tableFieldRes](
		H.MakeClient(http.DefaultClient),
	)
	HasField                = F.Curry2(uncurriedHasField)
	ResToReqTableWithParams = F.Curry2(uncurriedResToReqTableWithParams)
)

type tableFieldUpdateResponse struct {
	Id      int `json:"id"`
	TableId int `json:"table_id"`
}

type jsonPayload struct {
	Error   error
	Payload []byte
}

type tableFieldDelResponse struct {
	RelatedFields []struct {
		ID      int `json:"id"`
		TableID int `json:"table_id"`
	} `json:"related_fields,omitempty"`
}

type fieldsReqFeedback struct {
	Error  error
	Fields []tableFieldRes
	Msg    string
	Table  *client.Table
}

type tableFieldRes struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type tableFieldReq struct {
	tableFieldRes
	Params map[string]interface{}
}

func uncurriedHasField(name string, fieldResp tableFieldRes) bool {
	return fieldResp.Name == name
}

func onTableCreateFeedbackSuccess(res tableFieldRes) fieldsReqFeedback {
	return fieldsReqFeedback{
		Table: &client.Table{
			Id:   int32(res.Id),
			Name: res.Name,
		},
	}
}

func onFieldsReqFeedbackError(err error) fieldsReqFeedback {
	return fieldsReqFeedback{Error: err}
}

func onFieldsReqFeedbackSuccess(resp []tableFieldRes) fieldsReqFeedback {
	return fieldsReqFeedback{Fields: resp}
}

func onFieldDelReqFeedbackSuccess(
	resp tableFieldDelResponse,
) fieldsReqFeedback {
	return fieldsReqFeedback{Msg: "deleted field"}
}

func onFieldUpdateReqFeedbackSuccess(
	resp tableFieldUpdateResponse,
) fieldsReqFeedback {
	return fieldsReqFeedback{Msg: "updated field"}
}

func onFieldDelReqFeedbackNone() fieldsReqFeedback {
	return fieldsReqFeedback{Msg: "no field found to delete"}
}

func onJSONPayloadError(err error) jsonPayload {
	return jsonPayload{Error: err}
}

func onJSONPayloadSuccess(resp []byte) jsonPayload {
	return jsonPayload{Payload: resp}
}

func uncurriedResToReqTableWithParams(
	params map[string]interface{},
	req tableFieldRes,
) tableFieldReq {
	return tableFieldReq{
		tableFieldRes: tableFieldRes{
			Name: req.Name,
			Id:   req.Id,
		},
		Params: params,
	}
}

func ResToReqTable(req tableFieldRes) tableFieldReq {
	return tableFieldReq{
		tableFieldRes: tableFieldRes{
			Name: req.Name,
			Id:   req.Id,
		},
	}
}
