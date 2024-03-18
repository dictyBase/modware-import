package database

import (
	"context"
	"fmt"

	A "github.com/IBM/fp-go/array"
	O "github.com/IBM/fp-go/option"

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
	workspace, email string,
) (bool, int, error) {
	empty := 0
	ok := false
	wsp, err := wkm.ListWorkspaces()
	if err != nil {
		return ok, empty, err
	}
	output := F.Pipe3(
		wsp,
		A.FindFirst(HasWorkspace(workspace)),
		O.Map(SearchUser(email)),
		O.GetOrElse(F.Constant(0)),
	)
	if output != empty {
		ok = true
	}

	return ok, output, nil
}

func (wkm *WorkspaceManager) ListWorkspaces() ([]WorkspaceResp, error) {
	resp := F.Pipe3(
		wkm.ListWorkspaceURL(),
		H.MakeGetRequest,
		R.Map(httpapi.SetHeaderWithJWT(wkm.Token)),
		readWorkspaceResp,
	)(context.Background())
	output := F.Pipe1(
		resp(),
		E.Fold[error, []WorkspaceResp, listReqFeedback](
			onWorkspaceReqFeedbackError,
			onWorkspaceReqFeedbackSuccess,
		),
	)

	return output.Resp, output.Error
}
