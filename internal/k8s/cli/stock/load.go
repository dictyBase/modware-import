package stock

import (
	cliJob "github.com/dictyBase/modware-import/internal/k8s/cli/job"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var LoadStockCmd = &cobra.Command{
	Use:   "load-stock",
	Short: "run stock load command as a kubernetes job in the cluster",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		labels := cliJob.MetaLabel()
		labels["subcommand"] = "load-stock"
		for _, target := range viper.GetStringSlice("targets") {
			job, err := cliJob.Run(
				cmd,
				labels,
				LoadCommand(target),
			)
			if err != nil {
				return errors.Errorf("error in running job %s for target %s", err, target)
			}
			registry.GetLogger().Infof("deployed job %s for target %s", job.Name, target)
		}

		return nil
	},
}

func LoadCommand(target string) []string {
	return []string{
		"/usr/local/bin/gmake",
		target,
	}
}

func init() {
	LoadStockCmd.Flags().
		StringArray("targets", []string{"stock:loadStrain"}, "targets for loading stock data")
	viper.BindPFlags(LoadStockCmd.Flags())
}
