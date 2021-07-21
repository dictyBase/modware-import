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
	*k8s.Application
	LogLevel string
}

func NewUniprotLoader(args *k8s.AppParams, ispec *k8s.ImageSpec, level string) (*UniprotLoader, error) {
	app, err := k8s.NewApp(args, ispec)
	if err != nil {
		return &UniprotLoader{}, err
	}
	return &UniprotLoader{
		Application: app,
		LogLevel:    level,
	}, err
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

// Uninstall will attempt to uninstall in Kubernetes
func (u *UniprotLoader) Uninstall(client *kubernetes.Clientset) error {
	return client.BatchV1().
		Jobs(u.Meta().Namespace).
		Delete(
			context.Background(),
			u.Meta().Name,
			metav1.DeleteOptions{},
		)
}

func (u *UniprotLoader) Command() []string {
	return []string{
		"/usr/local/bin/app",
		"uniprot",
		"mapping",
	}
}

func (u *UniprotLoader) EnvManifest() []v1.EnvVar {
	return []v1.EnvVar{
		{
			Name: "ACCESS_KEY",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{Name: "dictybase-configuration"},
					Key:                  "minio.accesskey",
				},
			},
		},
		{
			Name: "SECRET_KEY",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{Name: "dictybase-configuration"},
					Key:                  "minio.secretkey",
				},
			},
		},
		{
			Name:  "LOG_LEVEL",
			Value: u.LogLevel,
		},
	}
}

func (u *UniprotLoader) TemplatePodSpecMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name: u.Meta().Name,
		Labels: map[string]string{
			"app": u.Meta().Name,
		},
	}
}

func (u *UniprotLoader) TemplatePodSpec() (apiv1.PodSpec, error) {
	contName, err := u.RandContainerName(10, "import")
	if err != nil {
		return apiv1.PodSpec{}, err
	}
	return apiv1.PodSpec{
		Containers: []apiv1.Container{
			{
				Name:    contName,
				Image:   u.ImageManifest(),
				Command: u.Command(),
				Env:     u.EnvManifest(),
			},
		},
	}, nil
}
