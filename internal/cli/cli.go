package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dictyBase/modware-import/internal/cli/stockcenter"
	"github.com/dictyBase/modware-import/internal/datasource/s3"
	"github.com/dictyBase/modware-import/internal/registry"
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
	RootCmd.AddCommand(stockcenter.StockCenterCmd)
	RootCmd.Flags().Bool("doc", false, "generate markdown documentation")
	RootCmd.PersistentFlags().String(
		"input-source",
		"bucket",
		"source of the file, could be one of bucket or folder",
	)
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
	RootCmd.PersistentFlags().StringP(
		"s3-server",
		"",
		"minio",
		"S3 server endpoint",
	)
	viper.BindEnv("s3-server", "MINIO_SERVICE_HOST")
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
	viper.BindPFlags(RootCmd.Flags())
	viper.BindPFlags(RootCmd.PersistentFlags())
}

// initConfig reads in config file and ENV variables if set.
//func initConfig() {
//if cfgFile != "" { // enable ability to specify config file via flag
//viper.SetConfigFile(cfgFile)
//}

//viper.SetConfigName(".modware-import") // name of config file (without extension)
//viper.AddConfigPath("$HOME")  // adding home directory as first search path
//viper.AutomaticEnv()          // read in environment variables that match

//// If a config file is found, read it in.
//if err := viper.ReadInConfig(); err == nil {
//fmt.Println("Using config file:", viper.ConfigFileUsed())
//}
//}
