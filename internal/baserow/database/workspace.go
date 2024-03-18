package database

import (
	"context"
	"fmt"

	R "github.com/IBM/fp-go/context/readerioeither"
	H "github.com/IBM/fp-go/context/readerioeither/http"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	"github.com/dictyBase/modware-import/internal/baserow/httpapi"

	"github.com/sirupsen/logrus"
)

type WorkspaceManager struct {
	Logger *logrus.Entry
	Token  string
	Host   string
}

func (wkm *WorkspaceManager) ListWorkspaceURL() string {
	return fmt.Sprintf("https://%s/api/workspaces/", wkm.Host)
}

func (wkm *WorkspaceManager) ListWorkspaceUserURL(id int, email string) string {
	return fmt.Sprintf(
		"https://%s/api/workspaces/users/workspace/%d?search=%s",
		wkm.Host, id, email,
	)
}

func (wkm *WorkspaceManager) SearchWorkspaceUser(
	workspaceId int,
	email string,
) (int, error) {
	resp := F.Pipe3(
		wkm.ListWorkspaceUserURL(workspaceId, email),
		H.MakeGetRequest,
		R.Map(httpapi.SetHeaderWithJWT(wkm.Token)),
		readWorkspaceUserResp,
	)(context.Background())
	output := F.Pipe1(
		resp(),
		E.Fold[error, []WorkspaceUserResp, workspaceUserFeedback](
			onWorkspaceUserReqFeedbackError,
			onWorkspaceUserReqFeedbackSuccess,
		),
	)

	return output.Resp[0].UserId, output.Error
}

func (wkm *WorkspaceManager) ListWorkspaces() ([]ListResponse, error) {
	resp := F.Pipe3(
		wkm.ListWorkspaceURL(),
		H.MakeGetRequest,
		R.Map(httpapi.SetHeaderWithJWT(wkm.Token)),
		readListResp,
	)(context.Background())
	output := F.Pipe1(
		resp(),
		E.Fold[error, []ListResponse, listReqFeedback](
			onListReqFeedbackError,
			onListReqFeedbackSuccess,
		),
	)

	return output.Resp, output.Error
}
