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

type DBCleaner struct {
	*k8s.SimpleJobApp
}

func NewDBCleaner(args *k8s.AppParams, ispec *k8s.ImageSpec, level string) (*DBCleaner, error) {
	app, err := k8s.NewSimpleJobApp(args, ispec, level)
	if err != nil {
		return &DBCleaner{}, err
	}
	return &DBCleaner{
		SimpleJobApp: app,
	}, err
}

func (d *DBCleaner) Command() []string {
	return []string{
		"/usr/local/bin/gmake",
		"cleanalldb",
	}
}

// Install will attempt to install in Kubernetes
func (d *DBCleaner) Install(client *kubernetes.Clientset) error {
	pspec, err := d.TemplatePodSpec()
	if err != nil {
		return err
	}
	job := &batch.Job{
		ObjectMeta: d.ObjectMeta,
		Spec: batch.JobSpec{
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: d.TemplatePodSpecMeta(),
				Spec:       pspec,
			},
		},
	}
	_, err = client.BatchV1().
		Jobs(d.Meta().Namespace).
		Create(
			context.Background(),
			job,
			metav1.CreateOptions{},
		)
	if err != nil {
		return errors.Errorf("error in creating dbcleaner job %s", err)
	}
	return nil
}

func (d *DBCleaner) TemplatePodSpec() (apiv1.PodSpec, error) {
	contName, err := k8s.RandContainerName(d.Meta().Name, "dbcleaner", 10)
	if err != nil {
		return apiv1.PodSpec{}, err
	}
	return apiv1.PodSpec{
		Containers: []apiv1.Container{
			{
				Name:    contName,
				Image:   d.ImageManifest(),
				Command: d.Command(),
				Env:     d.ContainerEnv(),
			},
		},
	}, nil
}

func (d *DBCleaner) ContainerEnv() []v1.EnvVar {
	v := append(
		k8s.ArangoConfigManifest(),
		k8s.LogEnvManifest(d.LogLevel())...,
	)
	return append(v, k8s.ArangoSecManifest()...)
}
