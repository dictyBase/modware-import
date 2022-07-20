package s3

import (
	"fmt"

	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/minio/minio-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var reqFlags = []string{"secret-key", "access-key", "s3-server", "s3-server-port"}

func NewS3Client(cmd *cobra.Command) (*minio.Client, error) {
	for _, flag := range reqFlags {
		if len(viper.GetString(flag)) == 0 {
			return &minio.Client{}, fmt.Errorf("required minio flag %s in missing", flag)
		}
	}
	registry.GetLogger().
		Debugf("using minio server:%s access-key:%s secret-key:%s port:%s",
			viper.GetString("s3-server"), viper.GetString("access-key"),
			viper.GetString("secret-key"), viper.GetString("s3-server-port"))
	return minio.New(
		fmt.Sprintf("%s:%s", viper.GetString("s3-server"), viper.GetString("s3-server-port")),
		viper.GetString("access-key"),
		viper.GetString("secret-key"),
		false,
	)
}
