package ontology

import (
	"github.com/cockroachdb/errors"
	cliJob "github.com/dictyBase/modware-import/internal/k8s/cli/job"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RefreshCmd = &cobra.Command{
	Use:   "refresh-ontology",
	Short: "run ontology refresh command as a kubernetes job in the cluster",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		labels := cliJob.MetaLabel()
		labels["subcommand"] = "refresh-ontology"
		job, err := cliJob.Run(cmd, labels, RefreshCommand())
		if err != nil {
			return errors.Errorf("error in running job %s", err)
		}
		registry.GetLogger().Infof("deployed job %s", job.Name)
		return nil
	},
}

func init() {
	RefreshCmd.Flags().String("branch", "master", "branch of github repository")
	RefreshCmd.Flags().String("group", "", "ontology group name")
	RefreshCmd.Flags().String("s3-bucket-path", "", "s3 bucket where files will be uploaded")
	_ = RefreshCmd.MarkFlagRequired("group")
	_ = RefreshCmd.MarkFlagRequired("s3-bucket-path")
	viper.BindPFlags(RefreshCmd.Flags())
}

func RefreshCommand() []string {
	return []string{
		"/usr/local/bin/importer",
		"ontology",
		"refresh",
		"--group",
		viper.GetString("group"),
		"--branch",
		viper.GetString("branch"),
		"--s3-bucket-path",
		viper.GetString("s3-bucket-path"),
	}
}
