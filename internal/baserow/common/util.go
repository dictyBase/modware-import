package common

import (
	"strings"
	"github.com/dictyBase/modware-import/internal/baserow/httpapi"
	J "github.com/IBM/fp-go/json"
	E "github.com/IBM/fp-go/either"
)

func ProcessOntologyTermId(val string) string {
	return strings.Replace(val, ":", "_", 1)
}

func MarshalPayload[T any](payload *T) E.Either[error, []byte] {
	return E.Pipe1(payload, J.Marshal)
}

func OnCreateFeedbackSuccess(annoId string, entityType string) httpapi.ResponseFeedback {
	return httpapi.ResponseFeedback{
		Msg: fmt.Sprintf("created %s with annotation id %s", entityType, annoId),
	}
}

func OnCreateFeedbackError(err error) httpapi.ResponseFeedback {
	return httpapi.ResponseFeedback{Err: err}
}
