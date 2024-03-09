package cli

import (
	"context"
	"fmt"

	"github.com/dictyBase/modware-import/internal/baserow/client"
	"github.com/dictyBase/modware-import/internal/baserow/database"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/urfave/cli/v2"
)

func CreateStrainTableHandler(cltx *cli.Context) error {
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
	strainTbl := &database.StrainTableManager{
		TableManager: &database.TableManager{
			Client:     database.BaserowClient(cltx.String("server")),
			DatabaseId: int32(cltx.Int("database-id")),
			Logger:     logger,
			Ctx:        authCtx,
			Token:      token,
		},
	}
	name := cltx.String("table")
	tbl, err := strainTbl.CreateTable(name, strainTbl.FieldNames())
	if err != nil {
		return cli.Exit(fmt.Sprintf("error in creating table %s", err), 2)
	}
	logger.Infof("created table with fields %s", tbl.GetName())
	flagNames := []string{
		"strainchar-ontology-table",
		"genetic-mod-ontology-table",
		"mutagenesis-method-ontology-table",
	}
	tableIdMaps, err := allTableIds(strainTbl.TableManager, flagNames, cltx)
	if err != nil {
		return cli.Exit(fmt.Sprintf("error in getting table ids %s", err), 2)
	}
	fieldDefs := []map[string]map[string]interface{}{
		strainTbl.LinkFieldChangeSpecs(tableIdMaps),
		strainTbl.FieldChangeSpecs(),
	}
	for _, def := range fieldDefs {
		err := updateFieldDefs(strainTbl.TableManager, def, tbl, logger)
		if err != nil {
			cli.Exit(err.Error(), 2)
		}
	}
	return nil
}