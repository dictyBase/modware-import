package database

import (
	"fmt"
	"net/http"

	A "github.com/IBM/fp-go/array"
	H "github.com/IBM/fp-go/context/readerioeither/http"
	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
	"github.com/dictyBase/modware-import/internal/baserow/client"
)

var (
	makeHTTPRequest  = F.Bind13of3(H.MakeRequest)
	readFieldDelResp = H.ReadJSON[tableFieldDelResponse](
		H.MakeClient(http.DefaultClient),
	)
	readFieldsResp = H.ReadJSON[[]tableFieldRes](
		H.MakeClient(http.DefaultClient),
	)
	readUpdateFieldsResp = H.ReadJSON[tableFieldUpdateResponse](
		H.MakeClient(http.DefaultClient),
	)
	readTableCreateResp = H.ReadJSON[tableFieldRes](
		H.MakeClient(http.DefaultClient),
	)
	readTablesResp = H.ReadJSON[[]tableFieldRes](
		H.MakeClient(http.DefaultClient),
	)
	readListRowsResp = H.ReadJSON[listRowsResp](
		H.MakeClient(http.DefaultClient),
	)
	readWorkspaceResp = H.ReadJSON[[]WorkspaceResp](
		H.MakeClient(http.DefaultClient),
	)
	HasField                   = F.Curry2(uncurriedHasField)
	ResToReqTableWithParams    = F.Curry2(uncurriedResToReqTableWithParams)
	matchTableName             = F.Curry2(uncurriedMatchTableName)
	onTablesReqFeedbackSuccess = F.Curry2(uncurriedOnTablesReqFeedbackSuccess)
	HasWorkspace               = F.Curry2(uncurriedHasWorkspace)
	HasUser                    = F.Curry2(uncurriedHasUser)
	SearchUser                 = F.Curry2(uncurriedSearchUser)
)

type rowResp struct {
	Id int32 `json:"id"`
}

type listRowsResp struct {
	Count    int32      `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []*rowResp `json:"results"`
}

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
	Id     int
}

type tableFieldRes struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type WorkspaceResp struct {
	Name  string              `json:"name"`
	Id    int                 `json:"id"`
	Users []WorkspaceUserResp `json:"users"`
}

type WorkspaceUserResp struct {
	Id    int    `json:"user_id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type listReqFeedback struct {
	Error error
	Msg   string
	Resp  []WorkspaceResp
}

type tableFieldReq struct {
	tableFieldRes
	Params map[string]interface{}
}

func uncurriedSearchUser(email string, wrsp WorkspaceResp) int {
	return F.Pipe3(
		wrsp.Users,
		A.FindFirst(HasUser(email)),
		O.Map(func(user WorkspaceUserResp) int { return user.Id }),
		O.GetOrElse(F.Constant(0)),
	)
}

func uncurriedHasField(name string, fieldResp tableFieldRes) bool {
	return fieldResp.Name == name
}

func uncurriedHasWorkspace(name string, wrsp WorkspaceResp) bool {
	return wrsp.Name == name
}

func uncurriedHasUser(email string, ursp WorkspaceUserResp) bool {
	return ursp.Email == email
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

func onListRowsReqFeedbackSuccess(resp listRowsResp) fieldsReqFeedback {
	return fieldsReqFeedback{Id: int(resp.Results[0].Id)}
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

func onWorkspaceReqFeedbackError(err error) listReqFeedback {
	return listReqFeedback{Error: err}
}

func onWorkspaceReqFeedbackSuccess(resp []WorkspaceResp) listReqFeedback {
	return listReqFeedback{Resp: resp}
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

func uncurriedMatchTableName(name string, tres tableFieldRes) bool {
	return tres.Name == name
}

func uncurriedOnTablesReqFeedbackSuccess(
	name string, res []tableFieldRes,
) fieldsReqFeedback {
	return F.Pipe2(
		res,
		A.FindFirst(matchTableName(name)),
		O.Fold(
			func() fieldsReqFeedback {
				return fieldsReqFeedback{
					Error: fmt.Errorf("unable to find table %s", name),
				}
			},
			func(ores tableFieldRes) fieldsReqFeedback {
				return fieldsReqFeedback{Id: ores.Id}
			},
		),
	)
}
