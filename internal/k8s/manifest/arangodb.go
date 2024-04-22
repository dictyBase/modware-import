package manifest

import (
	"fmt"

	apiv1 "k8s.io/api/core/v1"
)

func ArangoSecManifest(namespace string) []apiv1.EnvVar {
	return []apiv1.EnvVar{
		{
			Name: "ARANGODB_PASS",
			ValueFrom: &apiv1.EnvVarSource{
				SecretKeyRef: &apiv1.SecretKeySelector{
					LocalObjectReference: apiv1.LocalObjectReference{
						Name: fmt.Sprintf("dictycr-secret-%s", namespace),
					},
					Key: "arangodb.password",
				},
			},
		},
	}
}

func ArangoConfigManifest(namespace string) []apiv1.EnvVar {
	return []apiv1.EnvVar{
		{
			Name: "ARANGODB_USER",
			ValueFrom: &apiv1.EnvVarSource{
				SecretKeyRef: &apiv1.SecretKeySelector{
					LocalObjectReference: apiv1.LocalObjectReference{
						Name: fmt.Sprintf("dictycr-secret-%s", namespace),
					},
					Key: "arangodb.user",
				},
			},
		},
	}
}
