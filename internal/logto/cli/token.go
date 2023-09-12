package cli

import (
	"fmt"

	logto "github.com/dictyBase/modware-import/internal/logto/client"
	"github.com/urfave/cli/v2"
)

func CreateAccessToken(cltx *cli.Context) error {
	lclient := logto.NewClient(cltx.String("endpoint"))
	atoken, err := lclient.AccessToken(
		cltx.String("app-id"),
		cltx.String("app-secret"),
		cltx.String("api-resource"),
	)
	if err != nil {
		return cli.Exit(err, 2)
	}
	fmt.Printf("access token %s\n", atoken)
	return nil
}
