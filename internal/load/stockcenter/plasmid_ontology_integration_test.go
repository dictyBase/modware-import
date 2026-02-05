//go:build integration

package stockcenter

import (
	"context"
	"flag"
	"path/filepath"
	"strings"
	"testing"

	E "github.com/IBM/fp-go/either"
	"github.com/minio/minio-go/v6"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcminio "github.com/testcontainers/testcontainers-go/modules/minio"
	"github.com/urfave/cli/v2"
)

type s3TestFixture struct {
	container   testcontainers.Container
	minioClient *minio.Client
	endpoint    string
	accessKey   string
	secretKey   string
	cleanup     func()
}

func setupMinioContainer(t *testing.T) *s3TestFixture {
	t.Helper()
	ctx := context.Background()

	// Start MinIO testcontainer using official API
	minioContainer, err := tcminio.Run(ctx,
		"minio/minio:RELEASE.2024-01-16T16-07-38Z",
		tcminio.WithUsername("minioadmin"),
		tcminio.WithPassword("minioadmin"),
	)
	require.NoError(t, err, "failed to start MinIO container")

	// Get connection info
	endpoint, err := minioContainer.ConnectionString(ctx)
	require.NoError(t, err, "failed to get MinIO endpoint")

	// Create MinIO client
	minioClient, err := minio.New(
		endpoint,
		"minioadmin",
		"minioadmin",
		false, // no SSL
	)
	require.NoError(t, err, "failed to create MinIO client")

	return &s3TestFixture{
		container:   minioContainer,
		minioClient: minioClient,
		endpoint:    endpoint,
		accessKey:   "minioadmin",
		secretKey:   "minioadmin",
		cleanup: func() {
			if err := testcontainers.TerminateContainer(minioContainer); err != nil {
				t.Logf("failed to terminate container: %s", err)
			}
		},
	}
}

func (fix *s3TestFixture) uploadTestFile(
	t *testing.T,
	bucket, objectPath, localFilePath string,
) {
	t.Helper()

	// Create bucket if it doesn't exist
	exists, err := fix.minioClient.BucketExists(bucket)
	require.NoError(t, err, "failed to check bucket existence")

	if !exists {
		err = fix.minioClient.MakeBucket(bucket, "")
		require.NoError(t, err, "failed to create bucket")
	}

	// Upload file
	_, err = fix.minioClient.FPutObject(
		bucket,
		objectPath,
		localFilePath,
		minio.PutObjectOptions{},
	)
	require.NoError(t, err, "failed to upload file to MinIO")
}

func createS3TestContext(
	t *testing.T,
	fix *s3TestFixture,
	flags map[string]string,
) *cli.Context {
	t.Helper()

	set := flag.NewFlagSet("test", flag.ContinueOnError)

	// Required S3 flags
	set.String("s3-server", "", "s3 server")
	set.String("s3-server-port", "", "s3 port")
	set.String("access-key", "", "access key")
	set.String("secret-key", "", "secret key")
	set.String("s3-bucket", "", "bucket name")
	set.String("s3-bucket-path", "", "bucket path")
	set.String("input", "", "input file")
	set.String("input-source", "", "input source")

	// Set fixture values
	defaultFlags := map[string]string{
		"s3-server":      "localhost",
		"s3-server-port": extractPort(fix.endpoint),
		"access-key":     fix.accessKey,
		"secret-key":     fix.secretKey,
	}

	// Merge with provided flags
	for key, value := range defaultFlags {
		if _, exists := flags[key]; !exists {
			flags[key] = value
		}
	}

	for key, value := range flags {
		require.NoError(t, set.Set(key, value))
	}

	return cli.NewContext(nil, set, nil)
}

func extractPort(endpoint string) string {
	// Extract port from "localhost:XXXXX" format
	parts := strings.Split(endpoint, ":")
	if len(parts) == 2 {
		return parts[1]
	}
	return "9000"
}

func TestKeywordReaderFromS3BucketCli_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	fix := setupMinioContainer(t)
	defer fix.cleanup()

	// Upload test fixture file
	testFile := "plasmid_props_test.tsv"
	// Path is relative to project root (go up 3 levels: stockcenter -> load -> internal -> root)
	testFilePath := filepath.Join("..", "..", "..", "testdata", testFile)
	bucket := "test-bucket"
	s3Path := "import/data"

	fix.uploadTestFile(t, bucket, filepath.Join(s3Path, testFile), testFilePath)

	// Create CLI context
	cliCtx := createS3TestContext(t, fix, map[string]string{
		"s3-bucket":      bucket,
		"s3-bucket-path": s3Path,
		"input":          testFile,
		"input-source":   "bucket",
	})

	// Execute function
	config := KeywordLoaderCliConfig{
		Cmd: cliCtx,
	}

	result := keywordReaderFromS3BucketCli(config)()

	// Verify success
	require.True(t, E.IsRight(result), "expected Right, got Left")

	// Extract result
	cfg := E.GetOrElse(func(error) KeywordLoaderCliConfig { return config })(result)

	// Verify reader is configured
	require.NotNil(t, cfg.Reader, "CSV reader should be set")
	require.NotNil(t, cfg.Closer, "Closer should be set")

	// Verify TSV configuration
	require.Equal(t, '\t', cfg.Reader.Comma, "should use tab separator")
	require.Equal(t, -1, cfg.Reader.FieldsPerRecord, "should allow variable fields")

	// Read first line to verify content
	record, err := cfg.Reader.Read()
	require.NoError(t, err, "should read first record")
	require.Equal(t, 3, len(record), "should have 3 fields")
	require.Equal(t, "DBP0000034", record[0], "should read correct plasmid ID")
	require.Equal(t, "depositor", record[1], "should read correct property")
	require.Equal(t, "Gene Katz", record[2], "should read correct value")

	// Close the reader
	if cfg.Closer != nil {
		cfg.Closer.Close()
	}
}

func TestKeywordReaderFromS3BucketCli_InvalidCredentials(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	fix := setupMinioContainer(t)
	defer fix.cleanup()

	// Create context with INVALID credentials
	cliCtx := createS3TestContext(t, fix, map[string]string{
		"access-key":     "wrong-key",
		"secret-key":     "wrong-secret",
		"s3-bucket":      "test-bucket",
		"s3-bucket-path": "import/data",
		"input":          "plasmid_props_test.tsv",
	})

	config := KeywordLoaderCliConfig{Cmd: cliCtx}
	result := keywordReaderFromS3BucketCli(config)()

	// Verify error
	require.True(t, E.IsLeft(result), "expected error for invalid credentials")

	// Extract error and verify message
	leftResult := E.Swap(result)
	err := E.GetOrElse(func(KeywordLoaderCliConfig) error { return nil })(leftResult)
	require.NotNil(t, err, "should have error")
	require.Contains(
		t,
		err.Error(),
		"failed to stat s3 object",
		"error should mention S3 stat failure",
	)
}

func TestKeywordReaderFromS3BucketCli_FileNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	fix := setupMinioContainer(t)
	defer fix.cleanup()

	// Create bucket but DON'T upload file
	bucket := "test-bucket"
	err := fix.minioClient.MakeBucket(bucket, "")
	require.NoError(t, err, "failed to create bucket")

	cliCtx := createS3TestContext(t, fix, map[string]string{
		"s3-bucket":      bucket,
		"s3-bucket-path": "import/data",
		"input":          "nonexistent.tsv",
	})

	config := KeywordLoaderCliConfig{Cmd: cliCtx}
	result := keywordReaderFromS3BucketCli(config)()

	// Verify error
	require.True(t, E.IsLeft(result), "expected error for missing file")

	// Extract error and verify message
	leftResult := E.Swap(result)
	err = E.GetOrElse(func(KeywordLoaderCliConfig) error { return nil })(leftResult)
	require.NotNil(t, err, "should have error")
	require.Contains(
		t,
		err.Error(),
		"failed to stat s3 object",
		"error should mention S3 stat failure",
	)
}
