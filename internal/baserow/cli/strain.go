package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/dictyBase/modware-import/internal/baserow/client"
	"github.com/dictyBase/modware-import/internal/baserow/database"

	"github.com/dictyBase/modware-import/internal/baserow/strain"
	strainReader "github.com/dictyBase/modware-import/internal/datasource/xls/strain"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/urfave/cli/v2"
)

func flagNames() []string {
	allFlags := make([]string, 0)
	for _, flg := range strainOntologyTableFlags() {
		allFlags = append(allFlags, flg.Names()[0])
	}
	return allFlags
}

func LoadStrainAnnotationToTable(cltx *cli.Context) error {
	logger := registry.GetLogger()
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
	client := database.BaserowClient(cltx.String("server"))
	tbm := &database.TableManager{
		Client:     client,
		DatabaseId: int32(cltx.Int("database-id")),
		Logger:     logger,
		Ctx:        authCtx,
		Token:      token,
	}
	tableIdMaps, err := allTableIds(tbm, flagNames(), cltx)
	if err != nil {
		return cli.Exit(fmt.Sprintf("error in getting table ids %s", err), 2)
	}
	wkm := &database.WorkspaceManager{
		Logger: logger,
		Token:  token,
		Host:   client.GetConfig().Host,
	}
	loader := strain.NewStrainLoader(
		cltx.String("server"),
		token,
		cltx.String("workspace"),
		cltx.Int("table-id"),
		logger,
		tableIdMaps,
		tbm,
		wkm,
	)
	reader, err := strainReader.NewStrainAnnotationReader(
		cltx.String("input"),
		cltx.String("sheet"),
		time.Now(),
	)
	if err != nil {
		cli.Exit(err.Error(), 2)
	}
	if err := loader.Load(reader); err != nil {
		return cli.Exit(err.Error(), 2)
	}

	return nil
}

func CreateStrainTableHandler(cltx *cli.Context) error {
	logger := registry.GetLogger()
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
	tableIdMaps, err := allTableIds(strainTbl.TableManager, flagNames(), cltx)
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

func CreateStrainTableFlag() []cli.Flag {
	tblFlags := append(tableCreationFlags(),
		&cli.StringFlag{
			Name:     "table",
			Usage:    "table to create for loading strain annotation",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "strainchar-ontology-table",
			Usage: "table name containing strain characteristics ontology",
			Value: "strain_characteristics_ontology",
		},
		&cli.StringFlag{
			Name:  "genetic-mod-ontology-table",
			Usage: "table name containing genetic modification ontology",
			Value: "genetic_modification_ontology",
		},
		&cli.StringFlag{
			Name:  "mutagenesis-method-ontology-table",
			Usage: "table name containing mutagenesis method ontology",
			Value: "mutagenesis_method_ontology",
		},
	)
	return append(tblFlags, strainOntologyTableFlags()...)
}

func LoadStrainToTableFlag() []cli.Flag {
	tblFlags := append(tableCreationFlags(),
		&cli.StringFlag{
			Name:     "input",
			Aliases:  []string{"i"},
			Usage:    "input excel spreadsheet file with strain annotations",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "sheet",
			Aliases: []string{"s"},
			Usage:   "name of sheet which contains the annotation",
			Value:   "Strain_Annotations",
		},
		&cli.IntFlag{
			Name:     "table-id",
			Usage:    "Database table id",
			Required: true,
		},
	)
	return append(tblFlags, strainOntologyTableFlags()...)
}

func strainOntologyTableFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "strainchar-ontology-table",
			Usage: "table name containing strain characteristics ontology",
			Value: "strain_characteristics_ontology",
		},
		&cli.StringFlag{
			Name:  "genetic-mod-ontology-table",
			Usage: "table name containing genetic modification ontology",
			Value: "genetic_modification_ontology",
		},
		&cli.StringFlag{
			Name:  "mutagenesis-method-ontology-table",
			Usage: "table name containing mutagenesis method ontology",
			Value: "mutagenesis_method_ontology",
		},
	}
}
