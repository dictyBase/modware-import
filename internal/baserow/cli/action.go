package cli

import (
	"context"
	"fmt"

	"github.com/dictyBase/modware-import/internal/baserow/client"
	"github.com/dictyBase/modware-import/internal/baserow/database"
	"github.com/dictyBase/modware-import/internal/baserow/ontology"
	"github.com/dictyBase/modware-import/internal/collection"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slices"
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
	token, err := database.AccessToken(&database.AccessTokenProperties{
		Email:    cltx.String("email"),
		Password: cltx.String("password"),
		Server:   cltx.String("server"),
	})
	if err != nil {
		return cli.Exit(err, 2)
	}
	fmt.Println(token)
	return nil
}

func LoadOntologyToTable(cltx *cli.Context) error {
	logger := registry.GetLogger()
	bclient := database.BaserowClient(cltx.String("server"))
	authCtx := context.WithValue(
		context.Background(),
		client.ContextDatabaseToken,
		cltx.String("token"),
	)
	fields := map[string]client.Type712Enum{
		"Name":        client.TEXT,
		"Id":          client.TEXT,
		"Is_obsolete": client.BOOLEAN,
	}
	ok, err := database.CreateOntologyTableFields(
		&database.OntologyTableFieldsProperties{
			Client:   bclient,
			Logger:   logger,
			Ctx:      authCtx,
			FieldMap: fields,
			TableId:  cltx.Int("table-id"),
		},
	)
	if err != nil {
		return cli.Exit(err.Error(), 2)
	}
	props := &ontology.LoadProperties{
		File:    cltx.String("input"),
		TableId: cltx.Int("table-id"),
		Token:   cltx.String("token"),
		Client:  bclient,
		Logger:  logger,
	}
	if ok {
		if err := ontology.LoadNewOrUpdate(props); err != nil {
			return cli.Exit(err.Error(), 2)
		}
		return nil
	}
	if err := ontology.LoadNew(props); err != nil {
		return cli.Exit(err.Error(), 2)
	}

	return nil
}

func CreateTable(cltx *cli.Context) error {
	logger := registry.GetLogger()
	bclient := database.BaserowClient(cltx.String("server"))
	authCtx := context.WithValue(
		context.Background(),
		client.ContextAccessToken,
		cltx.String("token"),
	)
	tbl, resp, err := bclient.
		DatabaseTablesApi.
		CreateDatabaseTable(authCtx, int32(cltx.Int("database-id"))).
		TableCreate(client.TableCreate{Name: cltx.String("table")}).
		Execute()
	if err != nil {
		return cli.Exit(
			fmt.Errorf("error in creating table %s", err), 2,
		)
	}
	defer resp.Body.Close()
	logger.Infof("created table %s", tbl.GetName())
	return nil
}
