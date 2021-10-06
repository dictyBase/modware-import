package ontology

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dictyBase/go-genproto/dictybaseapis/api/upload"
	"github.com/dictyBase/modware-import/internal/cli/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	sreg "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/minio/minio-go/v6"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	chunkSize int = 1024
)

type streamOboParams struct {
	name    string
	reader  io.Reader
	gstream grpcStream
	logger  *logrus.Entry
}

type grpcStream interface {
	Send(*upload.FileUploadRequest) error
	CloseAndRecv() (*upload.FileUploadResponse, error)
	grpc.ClientStream
}

// LoadCmd obojson formmated ontologies to arangodb
var LoadCmd = &cobra.Command{
	Use:   "load",
	Short: "load an obojson formatted ontologies to arangodb",
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		for _, fn := range []func() error{stockcenter.SetAnnoAPIClient, stockcenter.SetStrainAPIClient} {
			if err := fn(); err != nil {
				return err
			}
		}
		return setOboReaders(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := registry.GetLogger()
		switch viper.GetString("grpc-client") {
		case sreg.ANNOTATION_CLIENT:
			return streamToAnnotationServer(logger)
		case sreg.StockClient:
			return streamToStockServer(logger)
		}
		logger.Infof(
			"uploaded %d obojson files",
			len(registry.GetAllReaders(registry.OboReadersKey)),
		)
		return nil
	},
}

func streamToAnnotationServer(logger *logrus.Entry) error {
	for name, reader := range registry.GetAllReaders(registry.OboReadersKey) {
		stream, err := sreg.GetAnnotationAPIClient().
			OboJSONFileUpload(context.Background())
		if err != nil {
			return err
		}
		err = streamOboToGrpcServer(&streamOboParams{
			name:    name,
			reader:  reader,
			gstream: stream,
			logger:  logger,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func streamToStockServer(logger *logrus.Entry) error {
	for name, reader := range registry.GetAllReaders(registry.OboReadersKey) {
		stream, err := sreg.GetStockAPIClient().
			OboJSONFileUpload(context.Background())
		if err != nil {
			return err
		}
		err = streamOboToGrpcServer(&streamOboParams{
			name:    name,
			reader:  reader,
			gstream: stream,
			logger:  logger,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func streamOboToGrpcServer(args *streamOboParams) error {
	buff := make([]byte, chunkSize)
	for {
		count, err := args.reader.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error in reading file %s", err)
		}
		err = args.gstream.Send(&upload.FileUploadRequest{Content: buff[:count]})
		if err != nil {
			return fmt.Errorf("error in streaming file content %s", err)
		}
	}
	ustatus, err := args.gstream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("error in closing file stream %s", err)
	}
	args.logger.Debugf(
		"upload %s obojson file with status %s",
		args.name, ustatus.Status.String(),
	)
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
	grpcFlags()
	viper.BindPFlags(LoadCmd.Flags())
}

func grpcFlags() {
	LoadCmd.Flags().String(
		"grpc-service-client",
		"c",
		"The grpc service client that will be used to upload obojson file, only annotation and stock are supported",
	)
	LoadCmd.Flags().String(
		"annotation-grpc-host",
		"annotation-api",
		"grpc host address for annotation service",
	)
	LoadCmd.Flags().String(
		"annotation-grpc-port",
		"",
		"grpc port for annotation service",
	)
	LoadCmd.Flags().String(
		"stock-grpc-host",
		"stock-api",
		"grpc host address for stock service",
	)
	LoadCmd.Flags().String(
		"stock-grpc-port",
		"",
		"grpc port for stock service",
	)
	viper.BindEnv("annotation-grpc-host", "ANNOTATION_API_SERVICE_HOST")
	viper.BindEnv("annotation-grpc-port", "ANNOTATION_API_SERVICE_PORT")
	viper.BindEnv("stock-grpc-port", "STOCK_API_SERVICE_PORT")
	viper.BindEnv("stock-grpc-host", "STOCK_API_SERVICE_HOST")
}
