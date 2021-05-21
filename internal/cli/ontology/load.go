package ontology

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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
		return setOboStorage(cmd)
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
				logger.Debugf("obograph %s does not exist, have to be loaded", k)
				err := ds.SaveOboGraphInfo(g)
				if err != nil {
					return fmt.Errorf("error in saving graph information %s", err)
				}
				nt, err := ds.SaveTerms(g)
				if err != nil {
					return fmt.Errorf("error in saving terms %s", err)
				}
				logger.Debugf("saved %d terms", nt)
				nr, err := ds.SaveRelationships(g)
				if err != nil {
					return fmt.Errorf("error in saving relationships %s", err)
				}
				logger.Debugf("saved %d relationships", nr)
				continue
			}
			logger.Debugf("obograph %s exist, have to be updated", k)
			if err := ds.UpdateOboGraphInfo(g); err != nil {
				return fmt.Errorf("error in updating graph information %s", err)
			}
			it, ut, err := ds.SaveOrUpdateTerms(g)
			if err != nil {
				return fmt.Errorf("error in updating terms %s", err)
			}
			logger.Debugf("saved: %d and updated: %d terms", it, ut)
			ur, err := ds.SaveNewRelationships(g)
			if err != nil {
				return fmt.Errorf("error in saving relationships %s", err)
			}
			logger.Debugf("updated %d relationships", ur)
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
		files, err := obojsonFiles(viper.GetString("folder"))
		if err != nil {
			return err
		}
		for _, v := range files {
			r, err := os.Open(v)
			if err != nil {
				return fmt.Errorf("error in opening file %s %s", v, err)
			}
			rds[filepath.Base(v)] = r
		}
	case stockcenter.BUCKET:
		doneCh := make(chan struct{})
		defer close(doneCh)
		for oinfo := range registry.GetS3Client().ListObjects(
			viper.GetString("s3-bucket"),
			viper.GetString("s3-bucket-path"),
			true, doneCh,
		) {
			sinfo, err := registry.GetS3Client().StatObject(
				viper.GetString("s3-bucket"), oinfo.Key,
				minio.StatObjectOptions{},
			)
			if err != nil {
				return fmt.Errorf(
					"error in getting information for object %s %s",
					oinfo.Key, err,
				)
			}
			tagOk := false
			var val string
		INNER:
			for t := range sinfo.UserMetadata {
				if "ontology-group" == strings.ToLower(t) {
					tagOk = true
					val = sinfo.UserMetadata[t]
					break INNER
				}
			}
			if !tagOk {
				registry.GetLogger().Warnf(
					"ontology-group metadata is not present for %s",
					sinfo.Key,
				)
				continue
			}
			if val != viper.GetString("group") {
				registry.GetLogger().Warnf(
					"ontology group metadata value %s did not match %s for %s",
					val, viper.GetString("group"), sinfo.Key,
				)
				continue
			}
			obj, err := registry.GetS3Client().GetObject(
				viper.GetString("s3-bucket"), sinfo.Key,
				minio.GetObjectOptions{},
			)
			if err != nil {
				return fmt.Errorf("error in getting object %s", oinfo.Key)
			}
			rds[sinfo.Key] = obj
		}
	default:
		return fmt.Errorf("error input source %s not supported", viper.GetString("input-source"))
	}
	registry.SetAllReaders(registry.OboReadersKey, rds)
	return nil
}

func setOboStorage(cmd *cobra.Command) error {
	tls, _ := cmd.Flags().GetBool("is-secure")
	cp := &araobo.ConnectParams{
		User:     viper.GetString("arangodb-user"),
		Pass:     viper.GetString("arangodb-pass"),
		Database: viper.GetString("arangodb-database"),
		Host:     viper.GetString("arangodb-host"),
		Port:     viper.GetInt("arangodb-port"),
		Istls:    tls,
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
	LoadCmd.Flags().StringP(
		"folder",
		"f",
		"",
		"input folder with obojson format files",
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
	LoadCmd.Flags().String(
		"group",
		"",
		"file belong to this ontology group will be uploaded. Only works for S3 storage",
	)
}
