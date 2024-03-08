package cli

import (
	"context"
	"fmt"
	"os"

	"slices"

	"github.com/dictyBase/modware-import/internal/baserow/client"
	"github.com/dictyBase/modware-import/internal/baserow/database"
	"github.com/dictyBase/modware-import/internal/collection"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/urfave/cli/v2"
)

func CreateDatabaseToken(cltx *cli.Context) error {
	atoken, err := database.AccessToken(&database.AccessTokenProperties{
		Email:    cltx.String("email"),
		Password: cltx.String("password"),
		Server:   cltx.String("server"),
	})
	if err != nil {
		return cli.Exit(fmt.Errorf("error in creating access token %s", err), 2)
	}
	bclient := database.BaserowClient(cltx.String("server"))
	authCtx := context.WithValue(
		context.Background(),
		client.ContextAccessToken,
		atoken,
	)
	wlist, r, err := bclient.WorkspacesApi.ListWorkspaces(authCtx).
		Execute()
	defer r.Body.Close()
	if err != nil {
		return cli.Exit(
			fmt.Errorf("error in executing list workspaces API call %s", err),
			2,
		)
	}
	wnames := collection.Map(
		wlist,
		func(w client.WorkspaceUserWorkspace) string { return w.GetName() },
	)
	idx := slices.Index(wnames, cltx.String("workspace"))
	if idx == -1 {
		return cli.Exit(
			fmt.Errorf(
				"workspace %s cannot be found",
				cltx.String("workspace"),
			),
			2,
		)
	}
	tok, r, err := bclient.DatabaseTokensApi.CreateDatabaseToken(authCtx).
		TokenCreate(client.TokenCreate{
			Name:      cltx.String("name"),
			Workspace: wlist[idx].GetId(),
		}).
		Execute()
	defer r.Body.Close()
	if err != nil {
		return cli.Exit(
			fmt.Errorf("error in creating token %s", err),
			2,
		)
	}
	fmt.Printf("database token %s\n", tok.GetKey())
	return nil
}

func CreateAccessToken(cltx *cli.Context) error {
	resp, err := database.AccessToken(&database.AccessTokenProperties{
		Email:    cltx.String("email"),
		Password: cltx.String("password"),
		Server:   cltx.String("server"),
	})
	if err != nil {
		return cli.Exit(err, 2)
	}
	fmt.Println(resp.GetToken())
	if cltx.Bool("save-refresh-token") {
		err := os.WriteFile(
			cltx.String("refresh-token-path"),
			[]byte(resp.GetRefreshToken()),
			0600,
		)
		if err != nil {
			return cli.Exit(
				fmt.Sprintf("error in writing refresh token to file %s", err),
				2,
			)
		}
		registry.GetLogger().
			Infof("saved refresh token at %s", cltx.String("refresh-token-path"))
	}
	return nil
}
