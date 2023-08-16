package manifest

import (
	apiv1 "k8s.io/api/core/v1"
)

func MinioEnv() []apiv1.EnvVar {
	return []apiv1.EnvVar{
		{
			Name: "ACCESS_KEY",
			ValueFrom: &apiv1.EnvVarSource{
				SecretKeyRef: &apiv1.SecretKeySelector{
					LocalObjectReference: apiv1.LocalObjectReference{Name: "dictycr-secret-dev"},
					Key:                  "minio.accesskey",
				},
			},
		},
		{
			Name: "SECRET_KEY",
			ValueFrom: &apiv1.EnvVarSource{
				SecretKeyRef: &apiv1.SecretKeySelector{
					LocalObjectReference: apiv1.LocalObjectReference{Name: "dictycr-sercet-dev"},
					Key:                  "minio.secretkey",
				},
			},
		},
	}
}
