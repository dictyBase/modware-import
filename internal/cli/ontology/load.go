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
		if err := setOboReaders(cmd); err != nil {
			return err
		}
		return setOboStorage(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := registry.GetLogger()
		ds := registry.GetArangoOboStorage()
		updated := 0
		fresh := 0
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
				logger.Debugf("saved %d terms", nt)
				nr, err := ds.SaveRelationships(g)
				if err != nil {
					return fmt.Errorf("error in saving relationships %s", err)
				}
				logger.Debugf("saved %d relationships", nr)
				fresh += 1
				continue
			}
			logger.Infof("obograph %s exist, have to be updated", k)
			if err := ds.UpdateOboGraphInfo(g); err != nil {
				return fmt.Errorf("error in updating graph information %s", err)
			}
			stats, err := ds.SaveOrUpdateTerms(g)
			if err != nil {
				return fmt.Errorf("error in updating terms %s", err)
			}
			logger.Debugf(
				"saved: %d updated: %d obsoleted: %d terms",
				stats.Created, stats.Updated, stats.Deleted,
			)
			ur, err := ds.SaveNewRelationships(g)
			if err != nil {
				return fmt.Errorf("error in saving relationships %s", err)
			}
			logger.Debugf("updated %d relationships", ur)
			updated += 1
		}
		logger.Infof(
			"loaded %d obo files, new %d updated %d",
			fresh+updated, fresh, updated,
		)
		return nil
	},
}

func setOboReaders(cmd *cobra.Command) error {
	switch viper.GetString("input-source") {
	case stockcenter.FOLDER:
		rds, err := setFileReaders()
		if err != nil {
			return err
		}
		registry.SetAllReaders(registry.OboReadersKey, rds)
	case stockcenter.BUCKET:
		rds, err := setBucketReaders(cmd)
		if err != nil {
			return err
		}
		registry.SetAllReaders(registry.OboReadersKey, rds)
	default:
		return fmt.Errorf(
			"error input source %s not supported",
			viper.GetString("input-source"),
		)
	}
	return nil
}

func setFileReaders() (map[string]io.Reader, error) {
	rds := make(map[string]io.Reader)
	files, err := obojsonFiles(viper.GetString("folder"))
	if err != nil {
		return rds, err
	}
	for _, v := range files {
		r, err := os.Open(v)
		if err != nil {
			return rds, fmt.Errorf("error in opening file %s %s", v, err)
		}
		rds[filepath.Base(v)] = r
	}
	return rds, nil
}

func setBucketReaders(cmd *cobra.Command) (map[string]io.Reader, error) {
	rds := make(map[string]io.Reader)
	doneCh := make(chan struct{})
	defer close(doneCh)
	for oinfo := range registry.GetS3Client().ListObjects(
		viper.GetString("s3-bucket"), viper.GetString("s3-bucket-path"),
		true, doneCh,
	) {
		sinfo, err := registry.GetS3Client().StatObject(
			viper.GetString("s3-bucket"), oinfo.Key, minio.StatObjectOptions{},
		)
		if err != nil {
			return rds, fmt.Errorf(
				"error in getting information for object %s %s", oinfo.Key, err,
			)
		}
		tagOk := false
		var val string
	INNER:
		for t := range sinfo.UserMetadata {
			if strings.ToLower(t) == GroupTag {
				tagOk = true
				val = sinfo.UserMetadata[t]
				break INNER
			}
		}
		if !tagOk {
			registry.GetLogger().Warnf(
				"ontology-group metadata is not present for %s", sinfo.Key,
			)
			continue
		}
		group, _ := cmd.Flags().GetString("group")
		if val != group {
			registry.GetLogger().Warnf(
				"ontology group metadata value %s did not match %s for %s",
				val, group, sinfo.Key,
			)
			continue
		}
		obj, err := registry.GetS3Client().GetObject(
			viper.GetString("s3-bucket"), sinfo.Key, minio.GetObjectOptions{},
		)
		if err != nil {
			return rds, fmt.Errorf("error in getting object %s", oinfo.Key)
		}
		rds[sinfo.Key] = obj
	}
	return rds, nil
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

func init() {
	LoadCmd.Flags().StringP(
		"folder",
		"f",
		"",
		"input folder with obojson format files",
	)
	LoadCmd.Flags().String(
		"group",
		"",
		"file belong to this ontology group will be uploaded. Only works for S3 storage",
	)
	LoadCmd.Flags().String(
		"grpc-service-client",
		"c",
		"The grpc service client that will be used to upload obojson file",
	)
	viper.BindPFlags(LoadCmd.Flags())
}
