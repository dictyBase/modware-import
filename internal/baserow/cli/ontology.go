package cli

import (
	"context"
	"fmt"
	"math"

	"github.com/dictyBase/modware-import/internal/baserow/client"
	"github.com/dictyBase/modware-import/internal/baserow/database"
	"github.com/dictyBase/modware-import/internal/baserow/ontology"
	"github.com/dictyBase/modware-import/internal/config"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/urfave/cli/v2"
)

func safeInt32(val int, flagName string) (int32, error) {
	if val > math.MaxInt32 || val < math.MinInt32 {
		return 0, fmt.Errorf("value %d for flag %s is out of int32 range", val, flagName)
	}
	return int32(val), nil
}

func LoadOntologyToTable(cltx *cli.Context) error {
	logger := registry.GetLogger()
	bclient := database.BaserowClient(cltx.String("server"))
	authCtx := context.WithValue(
		context.Background(),
		client.ContextDatabaseToken,
		cltx.String("token"),
	)
	dbID, err := safeInt32(cltx.Int("database-id"), "database-id")
	if err != nil {
		return cli.Exit(err.Error(), config.DefaultRetryBackoffFactor)
	}
	tblID, err := safeInt32(cltx.Int("table-id"), "table-id")
	if err != nil {
		return cli.Exit(err.Error(), config.DefaultRetryBackoffFactor)
	}

	ontTbl := &database.OntologyTableManager{
		TableManager: &database.TableManager{
			Client:     bclient,
			Logger:     logger,
			Ctx:        authCtx,
			Token:      cltx.String("token"),
			DatabaseId: dbID,
		},
	}
	ok, err := ontTbl.CheckAllTableFields(
		&client.Table{Id: tblID},
	)
	if err != nil {
		return cli.Exit(err.Error(), config.DefaultRetryBackoffFactor)
	}
	if !ok {
		return cli.Exit(
			fmt.Sprintf("table with id %d does not have the required fields", tblID),
			config.DefaultRetryBackoffFactor,
		)
	}
	props := &ontology.LoadProperties{
		File:    cltx.String("input"),
		TableId: cltx.Int("table-id"),
		Token:   cltx.String("token"),
		Client:  bclient,
		Logger:  logger,
	}
	if err := ontology.LoadNewOrUpdate(props); err != nil {
		return cli.Exit(err.Error(), config.DefaultRetryBackoffFactor)
	}

	return nil
}

func CreateOntologyTableHandler(cltx *cli.Context) error {
	token := cltx.String("token")
	if len(token) == 0 {
		rtoken, err := refreshToken(cltx)
		if err != nil {
			return cli.Exit(err.Error(), config.DefaultRetryBackoffFactor)
		}
		token = rtoken
	}
	authCtx := context.WithValue(
		context.Background(),
		client.ContextAccessToken,
		token,
	)
	dbID, err := safeInt32(cltx.Int("database-id"), "database-id")
	if err != nil {
		return cli.Exit(err.Error(), config.DefaultRetryBackoffFactor)
	}
	logger := registry.GetLogger()
	ontTbl := &database.OntologyTableManager{
		TableManager: &database.TableManager{
			Client:     database.BaserowClient(cltx.String("server")),
			Logger:     logger,
			Ctx:        authCtx,
			Token:      token,
			DatabaseId: dbID,
		},
	}
	for _, name := range cltx.StringSlice("table") {
		tbl, err := ontTbl.CreateTable(name, ontTbl.FieldNames())
		if err != nil {
			return cli.Exit(fmt.Sprintf("error in creating table %s", err), config.DefaultRetryBackoffFactor)
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
					config.DefaultRetryBackoffFactor,
				)
			}
			logger.Info(msg)
		}
	}
	return nil
}

func LoadOntologyToTableFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "token",
			Aliases:  []string{"t"},
			Usage:    "database token with write privilege",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "table-id",
			Usage:    "Database table id",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "input",
			Aliases:  []string{"i"},
			Usage:    "input json formatted ontology file",
			Required: true,
		},
	}
}

func CreateOntologyTableFlag() []cli.Flag {
	return append(tableCreationFlags(),
		&cli.StringSliceFlag{
			Name:     "table",
			Usage:    "tables to create for loading ontology",
			Required: true,
		},
	)
}
