package ontology

import (
	"github.com/cockroachdb/errors"
	cliJob "github.com/dictyBase/modware-import/internal/k8s/cli/job"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var LoadOntoCmd = &cobra.Command{
	Use:   "load-ontology",
	Short: "run ontology load command as a kubernetes job in the cluster",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		labels := cliJob.MetaLabel()
		labels["subcommand"] = "load-ontology"
		for _, dbname := range viper.GetStringSlice("databases") {
			job, err := cliJob.Run(cmd, labels, LoadCommand(dbname))
			if err != nil {
				return errors.Errorf("error in running job %s in database %s", err, dbname)
			}
			registry.GetLogger().Infof("deployed job %s in database %s", job.Name, dbname)
		}

		return nil
	},
}

func init() {
	LoadOntoCmd.Flags().
		StringArray("databases", []string{"stock", "annotation"}, "databases for loading ontologies")
	LoadOntoCmd.Flags().String("group", "", "ontology group name")
	LoadOntoCmd.Flags().String("s3-bucket-path", "", "s3 bucket from where files will be uploaded")
	_ = LoadOntoCmd.MarkFlagRequired("group")
	_ = LoadOntoCmd.MarkFlagRequired("s3-bucket-path")
	viper.BindPFlags(LoadOntoCmd.Flags())
}

func LoadCommand(dbname string) []string {
	return []string{
		"/usr/local/bin/importer",
		"ontology",
		"load",
		"--group",
		viper.GetString("group"),
		"--s3-bucket-path",
		viper.GetString("s3-bucket-path"),
		"--arangodb-database",
		dbname,
	}
}
