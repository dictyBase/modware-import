package httpapi

import (
	"fmt"
	"io"
	"net/http"

	H "github.com/IBM/fp-go/context/readerioeither/http"
	F "github.com/IBM/fp-go/function"
	B "github.com/IBM/fp-go/http/builder"
	C "github.com/IBM/fp-go/http/content"
	HD "github.com/IBM/fp-go/http/headers"
	S "github.com/IBM/fp-go/string"
)

var (
	WithJWT = F.Flow2(
		S.Format[string]("JWT %s"),
		B.WithAuthorization,
	)
	WithToken = F.Flow2(
		S.Format[string]("Token %s"),
		B.WithAuthorization,
	)
	MakeHTTPRequest = F.Bind13of3(H.MakeRequest)
)

func SetHeaderWithToken(token string) func(*http.Request) *http.Request {
	return func(req *http.Request) *http.Request {
		req.Header = F.Pipe3(
			B.Default,
			B.WithContentType(C.Json),
			B.WithHeader(HD.Accept)(C.Json),
			WithToken(token),
		).GetHeaders()

		return req
	}
}

func SetHeaderWith(req *http.Request) *http.Request {
	req.Header = F.Pipe2(
		B.Default,
		B.WithContentType(C.Json),
		B.WithHeader(HD.Accept)(C.Json),
	).GetHeaders()

	return req
}

func SetHeaderWithJWT(jwt string) func(*http.Request) *http.Request {
	return func(req *http.Request) *http.Request {
		req.Header = F.Pipe3(
			B.Default,
			B.WithContentType(C.Json),
			B.WithHeader(HD.Accept)(C.Json),
			WithJWT(jwt),
		).GetHeaders()

		return req
	}
}

func CommonHeader(lreq *http.Request, token, format string) {
	lreq.Header.Set("Content-Type", "application/json")
	lreq.Header.Set("Accept", "application/json")
	lreq.Header.Set("Authorization", fmt.Sprintf("%s %s", format, token))
}

func ReqToResponse(creq *http.Request) (*http.Response, error) {
	client := &http.Client{}
	uresp, err := client.Do(creq)
	if err != nil {
		return uresp, fmt.Errorf("error in making request %s", err)
	}
	if uresp.StatusCode != 200 {
		cnt, err := io.ReadAll(uresp.Body)
		if err != nil {
			return uresp, fmt.Errorf(
				"error in response and the reading the body %d %s",
				uresp.StatusCode,
				err,
			)
		}
		return uresp, fmt.Errorf(
			"unexpected error response %d %s",
			uresp.StatusCode,
			string(cnt),
		)
	}
	return uresp, nil
}
