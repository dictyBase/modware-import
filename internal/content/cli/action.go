package cli

import (
	"fmt"

	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/minio/minio-go/v6"
	"github.com/urfave/cli/v2"
)

func LoadContent(cltx *cli.Context) error {
	logger := registry.GetLogger()
	s3Client := registry.GetS3Client()
	// client := regsc.GetContentAPIClient()

	doneCh := make(chan struct{})
	defer close(doneCh)
	s3Objects := s3Client.ListObjects(
		cltx.String("s3-bucket"),
		cltx.String("s3-bucket-path"),
		true,
		doneCh,
	)
	for cinfo := range s3Objects {
		sinfo, err := s3Client.StatObject(
			cltx.String("s3-bucket"), cinfo.Key, minio.StatObjectOptions{},
		)
		if err != nil {
			return cli.Exit(
				fmt.Sprintf(
					"error in getting information for object %s %s",
					sinfo.Key, err,
				), 2)
		}
		logger.Infof("read content file %s", sinfo.Key)
		/* obj, err := s3Client.GetObject(
			cltx.String("s3-bucket"), sinfo.Key, minio.GetObjectOptions{},
		)
		if err != nil {
			return cli.Exit(fmt.Sprintf("error in gettting object %s", err), 2)
		} */
	}

	return nil
}
