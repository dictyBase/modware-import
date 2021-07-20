package app

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/dictyBase/modware-import/internal/k8s"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/apis/batch"
)

type UniprotLoader struct {
	*k8s.Application
}

func NewUniprotLoader(args *k8s.AppParams) (*UniprotLoader, error) {
	app, err := k8s.NewApp(args)
	if err != nil {
		return &UniprotLoader{}, err
	}
	return &UniprotLoader{Application: app}, err
}

// Install will attempt to install in Kubernetes
func (u *UniprotLoader) Install(client *kubernetes.Clientset) error {
	contName, err := u.RandContainerName(10, "import")
	if err != nil {
		return err
	}
	job := &batch.Job{
		ObjectMeta: u.Meta(),
		Spec: batch.JobSpec{
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: u.Meta().Name,
					Labels: map[string]string{
						"app": u.Meta().Name,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  contName,
							Image: "busybox",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	_, err = client.BatchV1().
		Jobs(u.Meta().Namespace).
		Create(context.Background(), job, metav1.CreateOptions{})
	if err != nil {
		return errors.Errorf("error in creating uniprot job %s", err)
	}
	return nil
}

// Uninstall will attempt to uninstall in Kubernetes
func (u *UniprotLoader) Uninstall(client *kubernetes.Clientset) error {
	panic("not implemented") // TODO: Implement
}
