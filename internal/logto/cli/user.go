package cli

import (
	logto "github.com/dictyBase/modware-import/internal/logto/client"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/urfave/cli/v2"
)

func ImportUser(cltx *cli.Context) error {
	logger := registry.GetLogger()
	lclient := logto.NewClient(cltx.String("endpoint"))
	aresp, err := lclient.AccessToken(
		cltx.String("app-id"),
		cltx.String("app-secret"),
		cltx.String("api-resource"),
	)
	if err != nil {
		return cli.Exit(err.Error(), 2)
	}
	userId, err := lclient.CreateUser(
		aresp.AccessToken,
		&logto.APIUsersPostReq{
			PrimaryEmail: "bola@bola.com",
			PrimaryPhone: "19343049303438",
			Username:     "hello",
			Password:     "r93r938493*7043",
			Name:         "bola",
		},
	)
	if err != nil {
		return cli.Exit(err.Error(), 2)
	}
	logger.Infof("got user id %s\n", userId)
	return nil
}
