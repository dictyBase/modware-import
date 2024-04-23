package ontology

import (
	"fmt"

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
		bucketPath, _ := cmd.Flags().GetString("s3-bucket-path-prefix")
		group, _ := cmd.Flags().GetString("group")
		for _, dbname := range viper.GetStringSlice("databases") {
			job, err := cliJob.Run(
				cmd,
				labels,
				LoadCommand(
					dbname,
					fmt.Sprintf("%s/%s", bucketPath, dbname),
					group,
				),
			)
			if err != nil {
				return errors.Errorf(
					"error in running job %s in database %s",
					err,
					dbname,
				)
			}
			registry.GetLogger().
				Infof("deployed job %s in database %s", job.Name, dbname)
		}

		return nil
	},
}

func init() {
	LoadOntoCmd.Flags().
		StringArray("databases", []string{"stock", "annotation"}, "databases for loading ontologies")
	LoadOntoCmd.Flags().String("group", "obojson", "ontology group name")
	LoadOntoCmd.Flags().
		String("s3-bucket-path-prefix", "import/obograph-json",
			"prefix of s3 bucket path from where the obojson files will be uploaded")
	viper.BindPFlags(LoadOntoCmd.Flags())
}

func LoadCommand(dbname, path, group string) []string {
	return []string{
		"/usr/local/bin/importer",
		"ontology",
		"load",
		"--log-level",
		viper.GetString("log-level"),
		"--group",
		group,
		"--s3-bucket-path",
		path,
		"--arangodb-database",
		dbname,
	}
}
