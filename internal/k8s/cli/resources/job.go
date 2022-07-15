package resources

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/spf13/cobra"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ListJobCmd = &cobra.Command{
	Use:   "list-jobs",
	Short: "prints list of jobs in a cluster under the given namespace",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		namespace, _ := cmd.Flags().GetString("namespace")
		client := registry.GetKubeClient(registry.KubeClientKey)
		jobList, err := client.BatchV1().
			Jobs(namespace).
			List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return errors.Errorf("error in retrieving job list %s", err)
		}
		for _, job := range jobList.Items {
			fmt.Printf("name: %s ==== status: %s\n", job.ObjectMeta.Name, status(job))
		}

		return nil
	},
}

func status(job batchv1.Job) string {
	switch {
	case job.Status.Succeeded > 0:
		return "success"
	case job.Status.Failed > 0:
		return "failed"
	}

	return "active"
}
