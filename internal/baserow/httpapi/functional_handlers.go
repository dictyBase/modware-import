package httpapi

type ResponseFeedback struct {
	Err error
	Msg string
}

type JSONPayload struct {
	Error   error
	Payload []byte
}

func OnJSONPayloadError(err error) JSONPayload {
	return JSONPayload{Error: err}
}

func OnJSONPayloadSuccess(resp []byte) JSONPayload {
	return JSONPayload{Payload: resp}
}
