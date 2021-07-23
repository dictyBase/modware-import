package app

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/dictyBase/modware-import/internal/k8s"
	batch "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type UniprotLoader struct {
	*k8s.SimpleJobApp
}

func NewUniprotLoader(args *k8s.AppParams, ispec *k8s.ImageSpec, level string) *UniprotLoader {
	return &UniprotLoader{
		SimpleJobApp: k8s.NewSimpleJobApp(args, ispec, level),
	}
}

func (u *UniprotLoader) Command() []string {
	return []string{
		"/usr/local/bin/importer",
		"uniprot",
		"mapping",
	}
}

// Install will attempt to install in Kubernetes
func (u *UniprotLoader) Install(client *kubernetes.Clientset) error {
	pspec, err := u.TemplatePodSpec()
	if err != nil {
		return err
	}
	job := &batch.Job{
		ObjectMeta: u.ObjectMeta,
		Spec: batch.JobSpec{
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: u.TemplatePodSpecMeta(),
				Spec:       pspec,
			},
		},
	}
	_, err = client.BatchV1().
		Jobs(u.Meta().Namespace).
		Create(
			context.Background(),
			job,
			metav1.CreateOptions{},
		)
	if err != nil {
		return errors.Errorf("error in creating uniprot job %s", err)
	}
	return nil
}

func (u *UniprotLoader) TemplatePodSpec() (apiv1.PodSpec, error) {
	contName, err := k8s.RandContainerName(u.Meta().Name, "import", 10)
	if err != nil {
		return apiv1.PodSpec{}, err
	}
	return apiv1.PodSpec{
		Containers: []apiv1.Container{
			{
				Name:    contName,
				Image:   u.ImageManifest(),
				Command: u.Command(),
				Env:     u.ContainerEnv(),
			},
		},
	}, nil
}

func (u *UniprotLoader) ContainerEnv() []v1.EnvVar {
	return append(
		k8s.MinioSecManifest(),
		k8s.LogEnvManifest(u.LogLevel())...,
	)
}
