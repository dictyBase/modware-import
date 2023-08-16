package manifest

import (
	apiv1 "k8s.io/api/core/v1"
)

func ArangoSecManifest() []apiv1.EnvVar {
	return []apiv1.EnvVar{
		{
			Name: "ARANGODB_PASS",
			ValueFrom: &apiv1.EnvVarSource{
				SecretKeyRef: &apiv1.SecretKeySelector{
					LocalObjectReference: apiv1.LocalObjectReference{
						Name: "dictycr-sercet-dev",
					},
					Key: "arangodb.password",
				},
			},
		},
	}
}

func ArangoConfigManifest() []apiv1.EnvVar {
	return []apiv1.EnvVar{
		{
			Name: "ARANGODB_USER",
			ValueFrom: &apiv1.EnvVarSource{
				SecretKeyRef: &apiv1.SecretKeySelector{
					LocalObjectReference: apiv1.LocalObjectReference{
						Name: "dictycr-sercet-dev",
					},
					Key: "arangodb.user",
				},
			},
		},
	}
}
