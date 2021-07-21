package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SimpleJobApp struct {
	metav1.ObjectMeta
	*ImageSpec
	description string
	randomizer  map[string]func() []string
}

func NewApp(args *AppParams, ispec *ImageSpec) (*SimpleJobApp, error) {
	qname, err := RandomFullName(args.Name, args.fragment, 10)
	if err != nil {
		return &SimpleJobApp{}, err
	}
	return &SimpleJobApp{
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
