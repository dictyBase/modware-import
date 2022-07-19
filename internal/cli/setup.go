package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/dictyBase/modware-import/internal/datasource/s3"
	"github.com/dictyBase/modware-import/internal/logger"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/minio/minio-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

func PersistentPreRun(cmd *cobra.Command) error {
	l, err := logger.NewLogger(cmd)
	if err != nil {
		return errors.Errorf("erron in getting a new logger %s", err)
	}
	registry.SetLogger(l)
	if len(viper.GetString("access-key")) > 0 && len(viper.GetString("secret-key")) > 0 {
		client, err := s3.NewS3Client(cmd)
		if err != nil {
			return errors.Errorf("error in getting instance of s3 client %s", err)
		}
		registry.SetS3Client(client)
	}
	return nil
}

func RunDoc(cmd *cobra.Command) error {
	mkdoc, _ := cmd.Flags().GetBool("doc")
	if mkdoc {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		docDir := filepath.Join(dir, "docs")
		if err := os.MkdirAll(docDir, 0700); err != nil {
			return err
		}
		if err := doc.GenMarkdownTree(cmd, docDir); err != nil {
			return err
		}
		fmt.Printf("created markdown docs in %s\n", docDir)
	}
	return nil
}

func PersistentPostRun(cmd *cobra.Command) error {
	if len(viper.GetString("access-key")) == 0 {
		return nil
	}
	if len(viper.GetString("secret-key")) == 0 {
		return nil
	}
	name := fmt.Sprintf("%s-%s.log", cmd.CalledAs(), time.Now().Format("20060102-150405"))
	_, err := registry.GetS3Client().FPutObject(
		viper.GetString("log-file-bucket"),
		fmt.Sprintf(
			"%s/%s",
			viper.GetString("log-file-bucket-path"),
			name,
		),
		registry.GetValue(registry.LogFileKey),
		minio.PutObjectOptions{},
	)
	if err != nil {
		return errors.Errorf(
			"error in uploading file %s with object name %s",
			registry.GetValue(registry.LogFileKey),
			name,
		)
	}
	return nil
}

func LoggingArgs(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(
		"log-level",
		"",
		"error",
		"log level for the application",
	)
	cmd.PersistentFlags().StringP(
		"log-format",
		"",
		"json",
		"format of the logging out, either of json or text",
	)
	cmd.PersistentFlags().String(
		"log-file",
		"",
		"file for log output other than standard output, written to a temp folder by default",
	)
	cmd.PersistentFlags().String(
		"log-file-bucket",
		"dictybase",
		"S3 bucket for log file",
	)
	cmd.PersistentFlags().String(
		"log-file-bucket-path",
		"import/log",
		"S3 path inside the bucket for storing log file",
	)
}

func S3Args(cmd *cobra.Command) {
	cmd.PersistentFlags().String(
		"s3-server-port",
		"",
		"S3 server port",
	)
	viper.BindEnv("s3-server-port", "MINIO_SERVICE_PORT")
	cmd.PersistentFlags().String(
		"access-key",
		"",
		"access key for S3 server",
	)
	viper.BindEnv("access-key", "ACCESS_KEY")
	cmd.PersistentFlags().String(
		"secret-key",
		"",
		"secret key for S3 server",
	)
	viper.BindEnv("access-key", "SECRET_KEY")
	cmd.PersistentFlags().String(
		"s3-server",
		"minio",
		"S3 server endpoint",
	)
	cmd.PersistentFlags().String(
		"s3-bucket",
		"dictybase",
		"S3 bucket for input files",
	)
	cmd.PersistentFlags().String(
		"s3-bucket-path",
		"",
		"path inside S3 bucket for input files[REQUIRED]",
	)
	_ = cmd.MarkFlagRequired("s3-bucket-path")
	viper.BindEnv("s3-server", "MINIO_SERVICE_HOST")
}
