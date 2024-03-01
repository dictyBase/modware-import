package cli

import (
	"context"
	"fmt"
	"os"

	"slices"

	"github.com/dictyBase/modware-import/internal/baserow/client"
	"github.com/dictyBase/modware-import/internal/baserow/database"
	"github.com/dictyBase/modware-import/internal/baserow/ontology"
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

func LoadOntologyToTable(cltx *cli.Context) error {
	logger := registry.GetLogger()
	bclient := database.BaserowClient(cltx.String("server"))
	authCtx := context.WithValue(
		context.Background(),
		client.ContextDatabaseToken,
		cltx.String("token"),
	)
	ontTbl := &database.OntologyTableManager{
		TableManager: &database.TableManager{
			Client:     bclient,
			Logger:     logger,
			Ctx:        authCtx,
			Token:      cltx.String("token"),
			DatabaseId: int32(cltx.Int("database-id")),
		},
	}
	ok, err := ontTbl.CheckAllTableFields(
		&client.Table{Id: int32(cltx.Int("table-id"))},
	)
	if err != nil {
		return cli.Exit(err.Error(), 2)
	}
	if !ok {
		return cli.Exit("table does not have the required fields", 2)
	}
	props := &ontology.LoadProperties{
		File:    cltx.String("input"),
		TableId: cltx.Int("table-id"),
		Token:   cltx.String("token"),
		Client:  bclient,
		Logger:  logger,
	}
	if err := ontology.LoadNewOrUpdate(props); err != nil {
		return cli.Exit(err.Error(), 2)
	}

	return nil
}

func CreatePhenoTableHandler(cltx *cli.Context) error {
	token := cltx.String("token")
	if len(token) == 0 {
		rtoken, err := refreshToken(cltx)
		if err != nil {
			return cli.Exit(err.Error(), 2)
		}
		token = rtoken
	}
	authCtx := context.WithValue(
		context.Background(),
		client.ContextAccessToken,
		token,
	)
	logger := registry.GetLogger()
	phenoTbl := &database.PhenotypeTableManager{
		TableManager: &database.TableManager{
			Client:     database.BaserowClient(cltx.String("server")),
			Logger:     logger,
			Ctx:        authCtx,
			Token:      token,
			DatabaseId: int32(cltx.Int("database-id")),
		},
	}
	name := cltx.String("table")
	tbl, err := phenoTbl.CreateTable(name, phenoTbl.FieldNames())
	if err != nil {
		return cli.Exit(fmt.Sprintf("error in creating table %s", err), 2)
	}
	logger.Infof("created table with fields %s", tbl.GetName())
	for fieldName, spec := range phenoTbl.FieldChangeSpecs() {
		msg, err := phenoTbl.UpdateField(tbl, fieldName, spec)
		if err != nil {
			return cli.Exit(
				fmt.Sprintf(
					"error in updating %s field %s",
					fieldName,
					err,
				),
				2,
			)
		}
		logger.Info(msg)
	}
	return nil
}

func CreateOntologyTableHandler(cltx *cli.Context) error {
	token := cltx.String("token")
	if len(token) == 0 {
		rtoken, err := refreshToken(cltx)
		if err != nil {
			return cli.Exit(err.Error(), 2)
		}
		token = rtoken
	}
	authCtx := context.WithValue(
		context.Background(),
		client.ContextAccessToken,
		token,
	)
	logger := registry.GetLogger()
	ontTbl := &database.OntologyTableManager{
		TableManager: &database.TableManager{
			Client:     database.BaserowClient(cltx.String("server")),
			Logger:     logger,
			Ctx:        authCtx,
			Token:      token,
			DatabaseId: int32(cltx.Int("database-id")),
		},
	}
	for _, name := range cltx.StringSlice("table") {
		tbl, err := ontTbl.CreateTable(name, ontTbl.FieldNames())
		if err != nil {
			return cli.Exit(fmt.Sprintf("error in creating table %s", err), 2)
		}
		logger.Infof("created table with fields %s", tbl.GetName())
		for fieldName, spec := range ontTbl.FieldChangeSpecs() {
			msg, err := ontTbl.UpdateField(tbl, fieldName, spec)
			if err != nil {
				return cli.Exit(
					fmt.Sprintf(
						"error in updating %s field %s",
						fieldName,
						err,
					),
					2,
				)
			}
			logger.Info(msg)
		}
	}
	return nil
}
func allPhenoTableIds(
	pheno *database.PhenotypeTableManager,
	cltx *cli.Context,
) (map[string]int, error) {
	idMaps := make(map[string]int)
	tableNames := []string{
		cltx.String("assay-ontology-table"),
		cltx.String("phenotype-ontology-table"),
		cltx.String("env-ontology-table"),
	}
	for _, name := range tableNames {
		id, err := pheno.TableNameToId(name)
		if err != nil {
			return idMaps, err
		}
		idMaps[name] = id
	}

	return idMaps, nil
}
