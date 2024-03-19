package cli

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	E "github.com/IBM/fp-go/either"

	A "github.com/IBM/fp-go/array"
	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
	S "github.com/IBM/fp-go/string"

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

func LoadStrainAnnotationFromFolderToTable(cltx *cli.Context) error {
	files, err := listStrainFiles(cltx.String("folder"))
	if err != nil {
		return cli.Exit(err.Error(), 2)
	}
	for _, rec := range files {
		if err := processFile(rec, cltx); err != nil {
			return cli.Exit(err.Error(), 2)
		}
	}
	return nil
}

func LoadStrainAnnotationToTable(cltx *cli.Context) error {
	err := processFile(cltx.String("input"), cltx)
	if err != nil {
		return cli.Exit(err.Error(), 2)
	}
	return nil
}

func processFile(filePath string, cltx *cli.Context) error {
	logger := registry.GetLogger()
	createdOn, err := parseStrainFileName(filePath)
	if err != nil {
		return err
	}
	reader, err := strainReader.NewStrainAnnotationReader(
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
	tableIdMaps, err := allTableIds(tbm, flagNames(), cltx)
	if err != nil {
		return fmt.Errorf("error in getting table ids %s", err)
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
	logger.Infof("going to load file %s", filePath)
	if err := loader.Load(reader); err != nil {
		return err
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
	)
	return append(tblFlags, strainOntologyTableFlags()...)
}

func LoadStrainFolderToTableFlag() []cli.Flag {
	return append(strainToTableFlag(),
		&cli.StringFlag{
			Name:     "folder",
			Aliases:  []string{"f"},
			Usage:    "folder with excel spreadsheet files with strain annotations",
			Required: true,
		})
}

func LoadStrainToTableFlag() []cli.Flag {
	return append(strainToTableFlag(),
		&cli.StringFlag{
			Name:     "input",
			Aliases:  []string{"i"},
			Usage:    "input excel spreadsheet file with strain annotations",
			Required: true,
		})
}

func strainToTableFlag() []cli.Flag {
	tblFlags := append(tableCreationFlags(),
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
		&cli.StringFlag{
			Name:     "workspace",
			Aliases:  []string{"w"},
			Usage:    "name of the workspace whether the database exists",
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

func parseStrainFileName(file string) (time.Time, error) {
	output := F.Pipe7(
		file,
		filepath.Base,
		Split("."),
		A.Head,
		O.GetOrElse(F.Constant("")),
		Split("-"),
		A.SliceRight[string](3),
		S.Join(":"),
	)
	if len(output) == 0 {
		return time.Time{}, fmt.Errorf("error in parsing file name %s", file)
	}
	return time.Parse("Jan:02:2006", output)
}

func listStrainFiles(folder string) ([]string, error) {
	output := F.Pipe2(
		E.TryCatchError(os.ReadDir(folder)),
		E.Map[error](func(files []fs.DirEntry) []string {
			return F.Pipe3(
				files,
				A.Filter(noDir),
				A.Filter(isStrainAnnoFile),
				A.Map(
					func(rec fs.DirEntry) string {
						return filepath.Join(folder, rec.Name())
					},
				),
			)
		}),
		E.Fold[error, []string](onErrorWithSlice, onSuccessWithSlice),
	)
	return output.Slice, output.Error
}

func isStrainAnnoFile(
	rec fs.DirEntry,
) bool {
	return F.Pipe1(rec.Name(), S.Includes("PMID"))
}
