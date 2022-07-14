package manifest

import (
	v1 "k8s.io/api/core/v1"
)

func MinioEnv() []v1.EnvVar {
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
	}
}
