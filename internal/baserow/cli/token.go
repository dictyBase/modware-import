package cli

import (
	"fmt"

	"github.com/dictyBase/modware-import/internal/baserow/httpapi"
	"github.com/urfave/cli/v2"
)

func refreshToken(cltx *cli.Context) (string, error) {
	var empty string
	tkm, err := httpapi.NewTokenManager(
		cltx.String("server"),
		cltx.String("refresh-token-path"),
	)
	if err != nil {
		return empty, err
	}
	token, err := tkm.FreshToken()
	if err != nil {
		return empty, fmt.Errorf("error in refreshing token %s", err)
	}

	return token, nil
}
