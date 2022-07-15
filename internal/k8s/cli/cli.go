package cli

import (
	"fmt"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/dictyBase/modware-import/internal/cli"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "k8s",
	Short: "cli for deploying and running import commands in kubernetes cluster",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.PersistentPreRun(cmd); err != nil {
			return errors.Errorf("error in executing pre run %s", err)
		}
		kcfg, _ := cmd.Flags().GetString("kubeconfig")
		if len(kcfg) == 0 {
			return nil
		}
		client, err := connectTok8s()
		if err != nil {
			return errors.Errorf("error in getting kubernetes client %s", err)
		}
		registry.SetKubeClient(registry.KubeClientKey, client)

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.RunDoc(cmd); err != nil {
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
	RootCmd.AddCommand(SubCmd)
	RootCmd.Flags().Bool("doc", false, "generate markdown documentation")
	RootCmd.PersistentFlags().String(
		"kubeconfig",
		"",
		"path to the kubernetes client(kubeconfig) file[REQUIRED]",
	)
	RootCmd.PersistentFlags().String("namespace", "dictybase", "kubernetes namespace")
	cli.LoggingArgs(RootCmd)
	cli.S3Args(RootCmd)
	viper.BindPFlags(RootCmd.Flags())
	viper.BindPFlags(RootCmd.PersistentFlags())
}

func connectTok8s() (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", viper.GetString("kubeconfig"))
	if err != nil {
		return &kubernetes.Clientset{}, errors.Errorf(
			"error in parsing config %s",
			err,
		)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return clientset, errors.Errorf("error in creating client from config %s", err)
	}

	return clientset, nil
}
