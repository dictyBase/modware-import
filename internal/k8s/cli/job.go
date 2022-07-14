package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// JobCmd is the base subcommand for grouping all commands that
// deploy kubernetes job manifest
var JobCmd = &cobra.Command{
	Use:   "job",
	Short: "subcommand for deploying kubernetes job manifest",
}

func init() {
	JobCmd.PersistentFlags().String("job", "", "name of the job")
	JobCmd.PersistentFlags().
		String("repo", "dictybase/modware-import", "container image repository")
	JobCmd.PersistentFlags().String("tag", "develop", "container image tag")
	_ = JobCmd.MarkFlagRequired("job")
	viper.BindPFlags(JobCmd.PersistentFlags())
}

func MetaLabel() map[string]string {
	return map[string]string{
		"command": "job",
		"runner":  "dictybot",
	}
}
