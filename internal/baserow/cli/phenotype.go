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
	tableIdMaps, err := allTableIds(phenoTbl.TableManager, flagNames, cltx)
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

func parsePhenoFileName(file string) (time.Time, error) {
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

func listPhenoFiles(folder string) ([]string, error) {
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

func isPhenoAnnoFile(
	rec fs.DirEntry,
) bool {
	return F.Pipe1(rec.Name(), S.Includes("PMID"))
}
