package strain

import (
	"fmt"
	"net/http"

	H "github.com/IBM/fp-go/context/readerioeither/http"
	"github.com/dictyBase/modware-import/internal/baserow/httpapi"
)

var (
	strainCreateHTTP = H.ReadJSON[strainCreateResp](
		H.MakeClient(http.DefaultClient),
	)
)

type strainCreateResp struct {
	AnnoId string `json:"annotation_id"`
}

func onStrainCreateFeedbackSuccess(
	res strainCreateResp,
) httpapi.ResponseFeedback {
	return httpapi.ResponseFeedback{
		Msg: fmt.Sprintf("created strain with annotation id %s", res.AnnoId),
	}
}

func onStrainCreateFeedbackError(err error) httpapi.ResponseFeedback {
	return httpapi.ResponseFeedback{Err: err}
}
