package cli

import (
	"fmt"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/dictyBase/modware-import/internal/cli/arangodb"
	"github.com/dictyBase/modware-import/internal/cli/data"
	"github.com/dictyBase/modware-import/internal/cli/ontology"
	"github.com/dictyBase/modware-import/internal/cli/stockcenter"
	"github.com/spf13/cobra"
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
		if err := PersistentPreRun(cmd); err != nil {
			return errors.Errorf("error in executing pre run %s", err)
		}
		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if err := PersistentPostRun(cmd); err != nil {
			return errors.Errorf("error in executing post-run %s", err)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := RunDoc(cmd); err != nil {
			return errors.Errorf("error in generating docs %s", err)
		}

		return nil
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.AddCommand(
		stockcenter.StockCenterCmd,
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
	LoggingArgs(RootCmd)
	S3Args(RootCmd)
	viper.BindPFlags(RootCmd.Flags())
	viper.BindPFlags(RootCmd.PersistentFlags())
}
