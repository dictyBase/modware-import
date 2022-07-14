package cli

import (
	"github.com/dictyBase/modware-import/internal/k8s/cli/ontology"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// JobCmd is the base subcommand for grouping all commands that
// deploy kubernetes job manifest
var SubCmd = &cobra.Command{
	Use:   "run",
	Short: "subcommand for deploying kubernetes manifest",
}

func init() {
	SubCmd.AddCommand(ontology.RefreshCmd, ontology.LoadOntoCmd)
	SubCmd.PersistentFlags().String("job", "", "name of the job")
	SubCmd.PersistentFlags().
		String("repo", "dictybase/modware-import", "container image repository")
	SubCmd.PersistentFlags().String("tag", "develop", "container image tag")
	_ = SubCmd.MarkFlagRequired("job")
	viper.BindPFlags(SubCmd.PersistentFlags())
}
