package cli

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	A "github.com/IBM/fp-go/array"
	Fn "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
	"github.com/dictyBase/go-genproto/dictybaseapis/content"
	"github.com/dictyBase/modware-import/internal/registry"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/minio/minio-go/v6"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var noncharReg = regexp.MustCompile("[^a-z0-9]+")

func Slugify(name string) string {
	return strings.Trim(
		noncharReg.ReplaceAllString(strings.ToLower(name), "-"),
		"-",
	)
}

func LoadContent(cltx *cli.Context) error {
	logger := registry.GetLogger()
	s3Client := registry.GetS3Client()
	client := regsc.GetContentAPIClient()

	doneCh := make(chan struct{})
	defer close(doneCh)
	s3Objects := listS3Objects(cltx, s3Client, doneCh)

	for cinfo := range s3Objects {
		err := processS3Object(cltx, logger, s3Client, client, cinfo)
		if err != nil {
			return cli.Exit(err.Error(), 2)
		}
	}

	return nil
}

func listS3Objects(
	cltx *cli.Context,
	s3Client *minio.Client,
	doneCh chan struct{},
) <-chan minio.ObjectInfo {
	return s3Client.ListObjects(
		cltx.String("s3-bucket"),
		cltx.String("s3-bucket-path"),
		true,
		doneCh,
	)
}

func processS3Object(
	cltx *cli.Context,
	logger *logrus.Entry,
	s3Client *minio.Client,
	client content.ContentServiceClient,
	cinfo minio.ObjectInfo,
) error {
	sinfo, err := s3Client.StatObject(
		cltx.String("s3-bucket"), cinfo.Key, minio.StatObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf(
			"error in getting information for object %s %s",
			sinfo.Key,
			err,
		)
	}
	logger.Infof("read content file %s", sinfo.Key)

	jsonContent, err := getContent(cltx, s3Client, sinfo)
	if err != nil {
		return err
	}

	name, namespace := nameAndNamespace(sinfo.Key)
	slug := Slugify(fmt.Sprintf("%s %s", name, namespace))

	err = storeOrUpdateContent(
		client,
		slug,
		name,
		namespace,
		jsonContent,
		logger,
		sinfo,
	)
	if err != nil {
		return err
	}

	return nil
}

func getContent(
	cltx *cli.Context,
	s3Client *minio.Client,
	sinfo minio.ObjectInfo,
) ([]byte, error) {
	obj, err := s3Client.GetObject(
		cltx.String("s3-bucket"), sinfo.Key, minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("error in getting object %s", err)
	}
	jsonContent, err := io.ReadAll(obj)
	if err != nil {
		return nil, fmt.Errorf(
			"error in reading content for file %s %s",
			sinfo.Key,
			err,
		)
	}

	return jsonContent, nil
}

func storeOrUpdateContent(
	client content.ContentServiceClient,
	slug, name, namespace string,
	jsonContent []byte,
	logger *logrus.Entry,
	sinfo minio.ObjectInfo,
) error {
	sct, err := client.GetContentBySlug(
		context.Background(),
		&content.ContentRequest{
			Slug: slug,
		},
	)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return createStoreContent(
				client,
				name,
				namespace,
				string(jsonContent),
				slug,
				logger,
				sinfo,
			)
		}
		return fmt.Errorf(
			"error in fetching content %s %s %s",
			sinfo.Key,
			sct.Data.Attributes.Name,
			sct.Data.Attributes.Namespace,
		)
	}
	logger.Infof("found existing content %s %s %s", sinfo.Key, name, namespace)

	return nil
}

func createStoreContent(
	client content.ContentServiceClient,
	name, namespace, jsonContent, slug string,
	logger *logrus.Entry,
	sinfo minio.ObjectInfo,
) error {
	nct, err := client.StoreContent(
		context.Background(),
		&content.StoreContentRequest{
			Data: &content.StoreContentRequest_Data{
				Attributes: &content.NewContentAttributes{
					Name:      name,
					Namespace: namespace,
					CreatedBy: "pfey@northwestern.edu",
					Content:   jsonContent,
					Slug:      slug,
				},
			},
		},
	)
	if err != nil {
		return fmt.Errorf(
			"error in creating content %s %s %s",
			sinfo.Key,
			name,
			namespace,
		)
	}
	logger.Infof(
		"created content %s %s %s",
		sinfo.Key,
		nct.Data.Attributes.Name,
		nct.Data.Attributes.Namespace,
	)

	return nil
}

func nameAndNamespace(input string) (string, string) {
	output := Fn.Pipe4(
		strings.Split(input, "/"),
		A.Last,
		O.Map(func(val string) []string { return strings.Split(val, ".") }),
		O.Map(func(val []string) string { return val[0] }),
		O.Map(func(val string) []string { return strings.Split(val, "-") }),
	)
	data, _ := O.Unwrap(output)
	return data[1], data[0]
}
