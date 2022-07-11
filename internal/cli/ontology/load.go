package ontology

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/dictyBase/go-obograph/graph"
	"github.com/dictyBase/go-obograph/storage"
	araobo "github.com/dictyBase/go-obograph/storage/arangodb"
	"github.com/dictyBase/modware-import/internal/cli/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/minio/minio-go/v6"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// LoadCmd obojson formmated ontologies to arangodb
var LoadCmd = &cobra.Command{
	Use:   "load",
	Short: "load obojson formatted ontologies to arangodb",
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return setOboReaders(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		dsa, err := araobo.NewDataSource(ConnectParams(), CollectParams())
		if err != nil {
			return errors.Errorf("error in connecting to arangodb %s", err)
		}
		logger := registry.GetLogger()
		for name, rdr := range registry.GetAllReaders(registry.OboReadersKey) {
			logger.Infof("going to load %s,", name)
			grph, err := graph.BuildGraph(rdr)
			if err != nil {
				return errors.Errorf("error in building graph from %s %s", name, err)
			}
			if !dsa.ExistsOboGraph(grph) {
				logger.Infof("obograph %s does not exist, have to be loaded", name)
				if err := saveNewGraph(dsa, grph, logger); err != nil {
					return errors.Errorf("error in saving new obograph %s %s", name, err)
				}

				continue
			}
			logger.Infof("obograph %s exist, have to be updated", name)
			if err := saveExistentGraph(dsa, grph, logger); err != nil {
				return errors.Errorf("error in saving existing obograph %s %s", name, err)
			}
		}
		logger.Infof(
			"uploaded %d obojson files",
			len(registry.GetAllReaders(registry.OboReadersKey)),
		)
		return nil
	},
}

func ConnectParams() *araobo.ConnectParams {
	arPort, _ := strconv.Atoi(viper.GetString("arangodb-port"))

	return &araobo.ConnectParams{
		User:     viper.GetString("arangodb-user"),
		Pass:     viper.GetString("arangodb-pass"),
		Host:     viper.GetString("arangodb-host"),
		Database: viper.GetString("arangodb-database"),
		Port:     arPort,
		Istls:    viper.GetBool("is-secure"),
	}
}

func CollectParams() *araobo.CollectionParams {
	return &araobo.CollectionParams{
		Term:         viper.GetString("term-collection"),
		Relationship: viper.GetString("rel-collection"),
		GraphInfo:    viper.GetString("cv-collection"),
		OboGraph:     viper.GetString("obograph"),
	}
}

func saveNewGraph(dsa storage.DataSource, grph graph.OboGraph, logger *logrus.Entry) error {
	err := dsa.SaveOboGraphInfo(grph)
	if err != nil {
		return errors.Errorf("error in saving graph %s", err)
	}
	nst, err := dsa.SaveTerms(grph)
	if err != nil {
		return errors.Errorf("error in saving terms %s", err)
	}
	logger.Infof("saved %d terms", nst)
	nsr, err := dsa.SaveRelationships(grph)
	if err != nil {
		return errors.Errorf("error in saving relationships %s", err)
	}
	logger.Infof("saved %d relationships", nsr)

	return nil
}

func saveExistentGraph(dsa storage.DataSource, grph graph.OboGraph, logger *logrus.Entry) error {
	if err := dsa.UpdateOboGraphInfo(grph); err != nil {
		return errors.Errorf("error in updating graph information %s", err)
	}
	stats, err := dsa.SaveOrUpdateTerms(grph)
	if err != nil {
		return errors.Errorf("error in updating terms %s", err)
	}
	logger.Infof(
		"saved::%d terms updated::%d terms obsoleted::%d terms",
		stats.Created, stats.Updated, stats.Deleted,
	)
	urs, err := dsa.SaveNewRelationships(grph)
	if err != nil {
		return errors.Errorf("error in saving relationships %s", err)
	}
	logger.Infof("updated %d relationships", urs)

	return nil
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
		return errors.Errorf(
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
			return rds, errors.Errorf("error in opening file %s %s", v, err)
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
			return rds, errors.Errorf(
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
			return rds, errors.Errorf("error in getting object %s", oinfo.Key)
		}
		rds[sinfo.Key] = obj
	}
	return rds, nil
}

func init() {
	LoadCmd.Flags().String(
		"term-collection",
		"cvterm",
		"arangodb collection for storing ontoloy terms",
	)
	LoadCmd.Flags().String(
		"rel-collection",
		"cvterm_relationship",
		"arangodb collection for storing cvterm relationships",
	)
	LoadCmd.Flags().String(
		"cv-collection",
		"cv",
		"arangodb collection for storing ontology information",
	)
	LoadCmd.Flags().String(
		"obograph",
		"obograph",
		"arangodb named graph for managing ontology graph",
	)
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
		"folder",
		"",
		"folder to read all obojson files",
	)
	viper.BindPFlags(LoadCmd.Flags())
}
