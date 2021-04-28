package ontology

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dictyBase/go-obograph/graph"
	araobo "github.com/dictyBase/go-obograph/storage/arangodb"
	"github.com/dictyBase/modware-import/internal/cli/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/minio/minio-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// LoadCmd obojson formmated ontologies to arangodb
var LoadCmd = &cobra.Command{
	Use:   "load",
	Short: "load an obojson formatted ontologies to arangodb",
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := setOboReaders(); err != nil {
			return err
		}
		return setOboStorage()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := registry.GetLogger()
		ds := registry.GetArangoOboStorage()
		for k, r := range registry.GetAllReaders(registry.OboReadersKey) {
			g, err := graph.BuildGraph(r)
			if err != nil {
				return err
			}
			if !ds.ExistsOboGraph(g) {
				logger.Infof("obograph %s does not exist, have to be loaded", k)
				err := ds.SaveOboGraphInfo(g)
				if err != nil {
					return fmt.Errorf("error in saving graph information %s", err)
				}
				nt, err := ds.SaveTerms(g)
				if err != nil {
					return fmt.Errorf("error in saving terms %s", err)
				}
				logger.Infof("saved %d terms", nt)
				nr, err := ds.SaveRelationships(g)
				if err != nil {
					return fmt.Errorf("error in saving relationships %s", err)
				}
				logger.Infof("saved %d relationships", nr)
				continue
			}
			logger.Infof("obograph %s exist, have to be updated", k)
			if err := ds.UpdateOboGraphInfo(g); err != nil {
				return fmt.Errorf("error in updating graph information %s", err)
			}
			it, ut, err := ds.SaveOrUpdateTerms(g)
			if err != nil {
				return fmt.Errorf("error in updating terms %s", err)
			}
			logger.Infof("saved: %d and updated: %d terms", it, ut)
			ur, err := ds.SaveNewRelationships(g)
			if err != nil {
				return fmt.Errorf("error in saving relationships %s", err)
			}
			logger.Infof("updated %d relationships", ur)
		}
		return nil
	},
}

func init() {
	loadFlags()
	viper.BindPFlags(LoadCmd.Flags())
}

func setOboReaders() error {
	rds := make(map[string]io.Reader)
	switch viper.GetString("input-source") {
	case stockcenter.FOLDER:
		for _, v := range viper.GetStringSlice("obojson") {
			r, err := os.Open(v)
			if err != nil {
				return fmt.Errorf("error in opening file %s %s", v, err)
			}
			rds[filepath.Base(v)] = r
		}
	case stockcenter.BUCKET:
		for _, v := range viper.GetStringSlice("obojson") {
			r, err := registry.GetS3Client().GetObject(
				viper.GetString("s3-bucket"),
				fmt.Sprintf("%s/%s", viper.GetString("s3-bucket-path"), v),
				minio.GetObjectOptions{},
			)
			if err != nil {
				return fmt.Errorf(
					"error in getting file %s from bucket %s %s",
					v, viper.GetString("s3-bucket-path"), err,
				)
			}
			rds[v] = r
		}
	default:
		return fmt.Errorf("error input source %s not supported", viper.GetString("input-source"))
	}
	registry.SetAllReaders(registry.OboReadersKey, rds)
	return nil
}

func setOboStorage() error {
	cp := &araobo.ConnectParams{
		User:     viper.GetString("arangodb-user"),
		Pass:     viper.GetString("arangodb-pass"),
		Database: viper.GetString("arangodb-database"),
		Host:     viper.GetString("arangodb-host"),
		Port:     viper.GetInt("arangodb-port"),
		Istls:    viper.GetBool("is-secure"),
	}
	clp := &araobo.CollectionParams{
		Term:         viper.GetString("term-collection"),
		Relationship: viper.GetString("rel-collection"),
		GraphInfo:    viper.GetString("cv-collection"),
		OboGraph:     viper.GetString("obograph"),
	}
	ds, err := araobo.NewDataSource(cp, clp)
	if err != nil {
		return err
	}
	registry.SetArangoOboStorage(ds)
	return nil
}

func loadFlags() {
	LoadCmd.Flags().StringSliceP(
		"obojson",
		"j",
		[]string{""},
		"input ontology files in obograph json format",
	)
	LoadCmd.Flags().String(
		"obograph",
		"obograph",
		"arangodb named graph for managing ontology graph",
	)
	LoadCmd.Flags().String(
		"cv-collection",
		"cv",
		"arangodb collection for storing ontology information",
	)
	LoadCmd.Flags().String(
		"rel-collection",
		"cvterm_relationship",
		"arangodb collection for storing cvterm relationships",
	)
	LoadCmd.Flags().String(
		"term-collection",
		"cvterm",
		"arangodb collection for storing ontoloy terms",
	)
}
