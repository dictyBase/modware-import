package cli

import (
	"context"
	"fmt"

	"github.com/dictyBase/modware-import/internal/baserow/client"
	"github.com/dictyBase/modware-import/internal/baserow/database"
	"github.com/dictyBase/modware-import/internal/baserow/ontology"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/urfave/cli/v2"
)

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
