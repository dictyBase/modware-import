package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dictyBase/modware-import/internal/cli/arangodb"
	"github.com/dictyBase/modware-import/internal/cli/data"
	"github.com/dictyBase/modware-import/internal/cli/ontology"
	"github.com/dictyBase/modware-import/internal/cli/stockcenter"
	"github.com/dictyBase/modware-import/internal/cli/uniprot"
	"github.com/dictyBase/modware-import/internal/datasource/s3"
	"github.com/dictyBase/modware-import/internal/logger"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/minio/minio-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "import",
	Short: "cli for importing dictybase data for migration",
	Long: `A command line application for importing dictyBase data.
The application is organized into subcommands which in turn has their
own subcommands for importing different kind of data. Each loading sub-subcommand
is generally expected to consume csv formatted data either directly from a source file
or through a file that is kept in a particular bucket of a S3 server.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if len(viper.GetString("access-key")) > 0 && len(viper.GetString("secret-key")) > 0 {
			client, err := s3.NewS3Client(cmd)
			if err != nil {
				return fmt.Errorf("error in getting instance of s3 client %s", err)
			}
			registry.SetS3Client(client)
		}
		l, err := logger.NewLogger(cmd)
		if err != nil {
			return err
		}
		registry.SetLogger(l)
		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
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
			return fmt.Errorf(
				"error in uploading file %s with object name %s",
				registry.GetValue(registry.LogFileKey),
				name,
			)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		d, _ := cmd.Flags().GetBool("doc")
		if d {
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
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.AddCommand(
		stockcenter.StockCenterCmd,
		uniprot.UniprotCmd,
		arangodb.ArangodbCmd,
		ontology.OntologyCmd,
		data.DataCmd,
	)
	RootCmd.Flags().Bool("doc", false, "generate markdown documentation")
	RootCmd.PersistentFlags().String(
		"input-source",
		"bucket",
		"source of the file, could be one of bucket or folder",
	)
	setLoggingArgs()
	setS3Args()
	viper.BindPFlags(RootCmd.Flags())
	viper.BindPFlags(RootCmd.PersistentFlags())
}

func setLoggingArgs() {
	RootCmd.PersistentFlags().StringP(
		"log-level",
		"",
		"error",
		"log level for the application",
	)
	RootCmd.PersistentFlags().StringP(
		"log-format",
		"",
		"json",
		"format of the logging out, either of json or text",
	)
	RootCmd.PersistentFlags().String(
		"log-file",
		"",
		"file for log output other than standard output, written to a temp folder by default",
	)
	RootCmd.PersistentFlags().String(
		"log-file-bucket",
		"dictybase",
		"S3 bucket for log file",
	)
	RootCmd.PersistentFlags().String(
		"log-file-bucket-path",
		"import/log",
		"S3 path inside the bucket for storing log file",
	)
}

func setS3Args() {
	RootCmd.PersistentFlags().StringP(
		"s3-server-port",
		"",
		"",
		"S3 server port",
	)
	viper.BindEnv("s3-server-port", "MINIO_SERVICE_PORT")
	RootCmd.PersistentFlags().StringP(
		"access-key",
		"",
		"",
		"access key for S3 server",
	)
	RootCmd.PersistentFlags().StringP(
		"secret-key",
		"",
		"",
		"secret key for S3 server",
	)
	RootCmd.PersistentFlags().StringP(
		"s3-server",
		"",
		"minio",
		"S3 server endpoint",
	)
	RootCmd.PersistentFlags().String(
		"s3-bucket",
		"dictybase",
		"S3 bucket for input files",
	)
	RootCmd.PersistentFlags().String(
		"s3-bucket-path",
		"",
		"path inside S3 bucket for input files",
	)
	viper.BindEnv("s3-server", "MINIO_SERVICE_HOST")
}
