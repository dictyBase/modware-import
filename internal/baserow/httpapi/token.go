package httpapi

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"

	R "github.com/IBM/fp-go/context/readerioeither"

	E "github.com/IBM/fp-go/either"
	J "github.com/IBM/fp-go/json"

	H "github.com/IBM/fp-go/context/readerioeither/http"
	F "github.com/IBM/fp-go/function"
)

var (
	readRefreshTokenResp = H.ReadJSON[refreshTokenRes](
		H.MakeClient(http.DefaultClient),
	)
)

type tokenReqFeedback struct {
	Error error
	Token string
}

type jsonPayload struct {
	Error   error
	Payload []byte
}

type TokenManager struct {
	host         string `validate:"required"`
	refreshToken string `validate:"required"`
}

type refreshTokenRes struct {
	Token string `json:"access_token"`
}

func NewTokenManager(
	host, tokenFile string,
) (*TokenManager, error) {
	cnt, err := os.ReadFile(tokenFile)
	if err != nil {
		return &TokenManager{}, fmt.Errorf(
			"error in reading from token file %s",
			err,
		)
	}
	tkm := &TokenManager{
		host:         host,
		refreshToken: string(cnt),
	}
	validate := validator.New()
	if err := validate.Struct(tkm); err != nil {
		return tkm, fmt.Errorf("validation failed %s", err)
	}

	return tkm, nil
}

func (tkm *TokenManager) RefreshTokenURL() string {
	return fmt.Sprintf(
		"https://%s/api/user/token-refresh/",
		tkm.host,
	)
}

func (tkm *TokenManager) FreshToken() (string, error) {
	var empty string
	tokenPayload := F.Pipe2(
		map[string]interface{}{
			"refresh_token": tkm.refreshToken,
		},
		J.Marshal,
		E.Fold(onJSONPayloadError, onJSONPayloadSuccess),
	)
	if tokenPayload.Error != nil {
		return empty, tokenPayload.Error
	}
	resp := F.Pipe3(
		tkm.RefreshTokenURL(),
		MakeHTTPRequest("POST", bytes.NewBuffer(tokenPayload.Payload)),
		R.Map(SetHeaderWith),
		readRefreshTokenResp,
	)(context.Background())
	output := F.Pipe1(
		resp(),
		E.Fold[error, refreshTokenRes, tokenReqFeedback](
			onTokenResError, onTokenResSuccess,
		),
	)

	return output.Token, output.Error
}

func onJSONPayloadError(err error) jsonPayload {
	return jsonPayload{Error: err}
}

func onJSONPayloadSuccess(resp []byte) jsonPayload {
	return jsonPayload{Payload: resp}
}

func onTokenResError(err error) tokenReqFeedback {
	return tokenReqFeedback{Error: err}
}

func onTokenResSuccess(res refreshTokenRes) tokenReqFeedback {
	return tokenReqFeedback{Token: res.Token}
}
