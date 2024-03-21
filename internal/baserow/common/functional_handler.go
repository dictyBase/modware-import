package common

import (
	"fmt"
	"net/http"
	"strings"

	H "github.com/IBM/fp-go/context/readerioeither/http"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	J "github.com/IBM/fp-go/json"
	"github.com/dictyBase/modware-import/internal/baserow/httpapi"
)

var (
	CreateHTTP = H.ReadJSON[CreateResp](H.MakeClient(http.DefaultClient))
)

type CreateResp struct {
	AnnoId string `json:"annotation_id"`
}

func ProcessOntologyTermId(val string) string {
	return strings.Replace(val, ":", "_", 1)
}

func MarshalPayload[T any](payload *T) E.Either[error, []byte] {
	return F.Pipe1(payload, J.Marshal)
}

func OnCreateFeedbackSuccess(
	annoId string,
	entityType string,
) httpapi.ResponseFeedback {
	return httpapi.ResponseFeedback{
		Msg: fmt.Sprintf(
			"created %s with annotation id %s",
			entityType,
			annoId,
		),
	}
}

func OnCreateFeedbackError(err error) httpapi.ResponseFeedback {
	return httpapi.ResponseFeedback{Err: err}
}
