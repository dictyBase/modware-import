package manifest

import (
	"k8s.io/api/core/v1"
)

func LogEnv(level string) []v1.EnvVar {
	return []v1.EnvVar{
		{
			Name:  "LOG_LEVEL",
			Value: level,
		},
	}
}
