package cli

import (
	"context"
	"fmt"

	"github.com/dictyBase/modware-import/internal/baserow/phenotype"
	phenoReader "github.com/dictyBase/modware-import/internal/datasource/xls/phenotype"

	"github.com/dictyBase/modware-import/internal/baserow/client"
	"github.com/dictyBase/modware-import/internal/baserow/database"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/urfave/cli/v2"
)

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
	flagNames := []string{
		"assay-ontology-table",
		"phenotype-ontology-table",
		"env-ontology-table",
	}
	tableIdMaps, err := allTableIDs(phenoTbl.TableManager, flagNames, cltx)
	if err != nil {
		return cli.Exit(fmt.Sprintf("error in getting table ids %s", err), 2)
	}
	fieldDefs := []map[string]map[string]interface{}{
		phenoTbl.LinkFieldChangeSpecs(tableIdMaps),
		phenoTbl.FieldChangeSpecs(),
	}
	for _, def := range fieldDefs {
		err := updateFieldDefs(phenoTbl.TableManager, def, tbl, logger)
		if err != nil {
			cli.Exit(err.Error(), 2)
		}
	}
	return nil
}

func LoadPhenoAnnotationToTable(cltx *cli.Context) error {
	err := processPhenoFile(cltx.String("input"), cltx)
	if err != nil {
		return cli.Exit(err.Error(), 2)
	}
	return nil
}

func LoadPhenoAnnotationFromFolderToTable(cltx *cli.Context) error {
	files, err := listPhenoFiles(cltx.String("folder"))
	if err != nil {
		return cli.Exit(err.Error(), 2)
	}
	for _, rec := range files {
		if err := processPhenoFile(rec, cltx); err != nil {
			return cli.Exit(err.Error(), 2)
		}
	}
	return nil
}

func processPhenoFile(filePath string, cltx *cli.Context) error {
	logger := registry.GetLogger()
	createdOn, err := parsePhenoFileName(filePath)
	if err != nil {
		return err
	}
	reader, err := phenoReader.NewPhenotypeAnnotationReader(
		filePath,
		cltx.String("sheet"),
		createdOn,
	)
	if err != nil {
		return err
	}
	token := cltx.String("token")
	if len(token) == 0 {
		token, err = refreshToken(cltx)
		if err != nil {
			return err
		}
	}
	authCtx := context.WithValue(
		context.Background(),
		client.ContextAccessToken,
		token,
	)
	client := database.BaserowClient(cltx.String("server"))
	tbm := &database.TableManager{
		Client:     client,
		DatabaseId: int32(cltx.Int("database-id")),
		Logger:     logger,
		Ctx:        authCtx,
		Token:      token,
	}
	tableIdMaps, err := allTableIDs(
		tbm,
		flagNamesHandler(phenoOntologyTableFlags()),
		cltx,
	)
	if err != nil {
		return fmt.Errorf("error in getting table ids %s", err)
	}
	wkm := &database.WorkspaceManager{
		Logger: logger,
		Token:  token,
		Host:   cltx.String("server"),
	}
	loader := phenotype.NewPhenotypeLoader(&phenotype.PhenotypeLoaderProperties{
		Host:             cltx.String("server"),
		Workspace:        cltx.String("workspace"),
		TableId:          cltx.Int("table-id"),
		Token:            token,
		Logger:           logger,
		OntologyTableMap: tableIdMaps,
		TableManager:     tbm,
		WorkspaceManager: wkm,
	})
	logger.Infof("going to load file %s", filePath)
	if err := loader.Load(reader); err != nil {
		return err
	}
	return nil
}
