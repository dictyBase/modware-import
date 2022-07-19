package manifest

import (
	v1 "k8s.io/api/core/v1"
)

func ArangoSecManifest() []v1.EnvVar {
	return []v1.EnvVar{
		{
			Name: "ARANGODB_PASS",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{Name: "dictybase-configuration"},
					Key:                  "arangodb.password",
				},
			},
		},
	}
}

func ArangoConfigManifest() []v1.EnvVar {
	return []v1.EnvVar{
		{
			Name: "ARANGODB_USER",
			ValueFrom: &v1.EnvVarSource{
				ConfigMapKeyRef: &v1.ConfigMapKeySelector{
					LocalObjectReference: v1.LocalObjectReference{Name: "dictybase-configuration"},
					Key:                  "arangodb.user",
				},
			},
		},
	}
}
