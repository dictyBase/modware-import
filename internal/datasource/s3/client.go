package s3

import (
	"fmt"

	minio "github.com/minio/minio-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewS3Client(cmd *cobra.Command) (*minio.Client, error) {
	return minio.New(
		fmt.Sprintf("%s:%s", viper.GetString("s3-server"), viper.GetString("s3-server-port")),
		viper.GetString("access-key"),
		viper.GetString("secret-key"),
		false,
	)
}
