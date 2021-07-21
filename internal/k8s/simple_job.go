package k8s

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type SimpleJobApp struct {
	metav1.ObjectMeta
	*ImageSpec
	description string
	level       string
	randomizer  map[string]func() []string
}

func NewApp(args *AppParams, ispec *ImageSpec, level string) (*SimpleJobApp, error) {
	qname, err := RandomFullName(args.Name, args.fragment, 10)
	if err != nil {
		return &SimpleJobApp{}, err
	}
	return &SimpleJobApp{
		level:       level,
		description: args.Description,
		ImageSpec:   ispec,
		ObjectMeta: metav1.ObjectMeta{
			Name:            qname,
			Namespace:       args.Namespace,
			ResourceVersion: "v1.0.0",
			Labels: map[string]string{
				"heritage": "naml",
			},
		},
	}, nil
}

// Meta returns the Kubernetes native ObjectMeta which is used to manage applications with naml.
func (a *SimpleJobApp) Meta() *metav1.ObjectMeta {
	return &a.ObjectMeta
}

// Description returns the application description
func (a *SimpleJobApp) Description() string {
	return a.description
}

func (a *SimpleJobApp) TemplatePodSpecMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name: a.Meta().Name,
		Labels: map[string]string{
			"app": a.Meta().Name,
		},
	}
}

// Uninstall will attempt to uninstall in Kubernetes
func (a *SimpleJobApp) Uninstall(client *kubernetes.Clientset) error {
	return client.BatchV1().
		Jobs(a.Meta().Namespace).
		Delete(
			context.Background(),
			a.Meta().Name,
			metav1.DeleteOptions{},
		)
}

func (a *SimpleJobApp) LogLevel() string {
	return a.level
}
