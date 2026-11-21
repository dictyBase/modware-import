//go:build integration

package gbexport_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcminio "github.com/testcontainers/testcontainers-go/modules/minio"
	"gocloud.dev/blob"
	"gocloud.dev/blob/s3blob"

	"github.com/dictyBase/modware-import/internal/goldenbraid/gbexport"
)

const (
	minioUser   = "minioadmin"
	minioSecret = "minioadmin"
	minioRegion = "us-east-1"
)

type minioFixture struct {
	endpoint string
	client   *s3.Client
}

func startMinio(t *testing.T) minioFixture {
	t.Helper()
	ctx := context.Background()
	container, err := tcminio.Run(
		ctx,
		"minio/minio:RELEASE.2024-01-16T16-07-38Z",
		tcminio.WithUsername(minioUser),
		tcminio.WithPassword(minioSecret),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = testcontainers.TerminateContainer(container) })
	address, err := container.ConnectionString(ctx)
	require.NoError(t, err)
	endpoint := "http://" + address
	client := s3.New(s3.Options{
		BaseEndpoint: aws.String(endpoint),
		Region:       minioRegion,
		UsePathStyle: true,
		Credentials:  credentials.NewStaticCredentialsProvider(minioUser, minioSecret, ""),
	})
	return minioFixture{endpoint: endpoint, client: client}
}

func (m minioFixture) bucket(t *testing.T, name string) *blob.Bucket {
	t.Helper()
	ctx := context.Background()
	_, err := m.client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(name)})
	require.NoError(t, err)
	bucket, err := s3blob.OpenBucket(ctx, m.client, name, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = bucket.Close() })
	return bucket
}

func (m minioFixture) params(input, output string) gbexport.Params {
	return gbexport.Params{
		Input:      input,
		Output:     output,
		S3Endpoint: m.endpoint,
		S3Region:   minioRegion,
		AccessKey:  minioUser,
		SecretKey:  minioSecret,
	}
}

func requireObjects(t *testing.T, bucket *blob.Bucket, prefix string, stems ...string) {
	t.Helper()
	for _, stem := range stems {
		for _, key := range []string{
			fmt.Sprintf("%s/%s/%s.fasta", prefix, gbexport.FastaDir, stem),
			fmt.Sprintf("%s/%s/%s.gbk", prefix, gbexport.GenBankDir, stem),
			fmt.Sprintf("%s/%s/%s.svg", prefix, gbexport.SVGDir, stem),
		} {
			data, err := bucket.ReadAll(context.Background(), key)
			require.NoError(t, err, key)
			require.NotEmpty(t, data, key)
		}
	}
}

func TestS3ArchiveToS3Output(t *testing.T) {
	m := startMinio(t)
	bucket := m.bucket(t, "zipcase")
	require.NoError(
		t,
		bucket.WriteAll(context.Background(), "prefix/assemblies.zip", assemblyArchive(t), nil),
	)
	logs := &capture{}
	_, err := run(
		context.Background(),
		t,
		logs,
		m.params("s3://zipcase/prefix/assemblies.zip", "s3://zipcase/out"),
	)
	require.NoError(t, err)
	requireObjects(t, bucket, "out", allStems...)
}

func TestS3PrefixPartialFailure(t *testing.T) {
	m := startMinio(t)
	bucket := m.bucket(t, "dircase")
	ctx := context.Background()
	seed := map[string][]byte{
		"prefix/" + simpleStem + ".dna":                       fixture(t, simpleStem),
		"prefix/GoldenBraid Assemblies/" + denseStem + ".dna": fixture(t, denseStem),
		"prefix/corrupt.dna":                                  []byte("not a snapgene file"),
		"prefix/a/b/" + gapStem + ".dna":                      fixture(t, gapStem),
		"prefixb/" + gapStem + ".dna":                         fixture(t, gapStem),
	}
	for key, data := range seed {
		require.NoError(t, bucket.WriteAll(ctx, key, data, nil))
	}
	logs := &capture{}
	_, err := run(ctx, t, logs, m.params("s3://dircase/prefix", "s3://dircase/out"))
	require.ErrorIs(t, err, gbexport.ErrPartialExportFailure)
	requireObjects(t, bucket, "out", simpleStem, denseStem)
	requireSummary(t, logs, 3, 2, 1)
	failures := logs.withMessage(t, "export failed")
	require.Len(t, failures, 1)
	require.Equal(t, "prefix/corrupt.dna", failures[0]["key"])
}

func TestS3ArchiveToLocalOutput(t *testing.T) {
	m := startMinio(t)
	bucket := m.bucket(t, "mixed")
	require.NoError(
		t,
		bucket.WriteAll(context.Background(), "prefix/assemblies.zip", assemblyArchive(t), nil),
	)
	out := t.TempDir()
	logs := &capture{}
	_, err := run(context.Background(), t, logs, m.params("s3://mixed/prefix/assemblies.zip", out))
	require.NoError(t, err)
	requireArtifacts(t, out, allStems...)
	_, err = os.Stat(filepath.Join(out, gbexport.SVGDir))
	require.NoError(t, err)
}
