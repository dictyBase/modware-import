package env

import (
	"fmt"
	"os"
)

func ArangoEnvs() error {
	return checkErrors([]string{
		"ARANGODB_PASS",
		"ARANGODB_USER",
		"ARANGODB_SERVICE_HOST",
		"ARANGODB_SERVICE_PORT",
	})
}

func ServiceEnvs() error {
	return checkErrors([]string{
		"STOCK_API_SERVICE_HOST",
		"STOCK_API_SERVICE_PORT",
		"ANNOTATION_API_SERVICE_HOST",
		"ANNOTATION_API_SERVICE_PORT",
	})
}

func MinioEnvs() error {
	return checkErrors([]string{
		"ACCESS_KEY",
		"SECRET_KEY",
	})
}

func MinioAccessKey() string {
	return os.Getenv("ACCESS_KEY")
}

func MinioSecretKey() string {
	return os.Getenv("SECRET_KEY")
}

func checkErrors(envs []string) error {
	for _, e := range envs {
		v := os.Getenv(e)
		if len(v) == 0 {
			return fmt.Errorf("env %s is not set", e)
		}
	}
	return nil
}
