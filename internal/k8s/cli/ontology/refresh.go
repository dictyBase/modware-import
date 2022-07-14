package ontology

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/dictyBase/modware-import/internal/k8s/cli"
	"github.com/dictyBase/modware-import/internal/k8s/manifest"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var RefreshCmd = &cobra.Command{
	Use:   "refresh-ontology",
	Short: "run ontology refresh command as a kubernetes job in the cluster",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		labels := cli.MetaLabel()
		labels["subcommand"] = "refresh-ontology"
		jobManifest, err := manifest.NewJob(&manifest.JobParams{
			Cli:        cmd,
			Labels:     labels,
			Command:    RefreshCommand(),
			Fragment:   cli.Fragment,
			NameLength: cli.NameLen,
		}).MakeSpec()
		if err != nil {
			return errors.Errorf("error in making job manifest %s", err)
		}
		job, err := registry.GetKubeClient(registry.KubeClientKey).BatchV1().
			Jobs(viper.GetString("namespace")).
			Create(context.Background(), jobManifest, metav1.CreateOptions{})
		if err != nil {
			return errors.Errorf("error in deploying job %s", err)
		}
		registry.GetLogger().Infof("deployed job %s", job.Name)
		return nil
	},
}

func init() {
	RefreshCmd.Flags().String("branch", "master", "branch of github repository")
	RefreshCmd.Flags().String("group", "", "ontology group name")
	_ = RefreshCmd.MarkFlagRequired("group")
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
	}
}
