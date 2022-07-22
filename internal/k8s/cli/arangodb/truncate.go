package arangodb

import (
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cockroachdb/errors"
	cliJob "github.com/dictyBase/modware-import/internal/k8s/cli/job"
)

var TruncateCmd = &cobra.Command{
	Use:   "truncate-arangodb",
	Short: "run arangodb delete command as a kubernetes job in the cluster",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		labels := cliJob.MetaLabel()
		labels["subcommand"] = "truncate-arangodb"
		database, _ := cmd.Flags().GetString("database")
		job, err := cliJob.Run(
			cmd,
			labels,
			LoadCommand(database),
		)
		if err != nil {
			return errors.Errorf("error in running job %s for database %s", err, database)
		}
		registry.GetLogger().Infof("deployed job %s for database %s", job.Name, database)

		return nil
	},
}

func init() {
	TruncateCmd.Flags().
		String("database", "",
			"name of arangodb database whose data from all the non-system collections will be deleted")
	_ = TruncateCmd.MarkFlagRequired("database")
	viper.BindPFlags(TruncateCmd.Flags())
}

func LoadCommand(database string) []string {
	return []string{
		"/usr/local/bin/importer",
		"arangodb",
		"delete",
		"--log-level",
		"info",
		"--database",
		database,
		"--is-secure",
	}
}
