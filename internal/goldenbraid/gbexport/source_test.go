package gbexport

import (
	"context"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/require"
)

func TestS3OptionsTargetPathStyleEndpoint(t *testing.T) {
	opts := s3Options(Params{
		S3Endpoint: "http://localhost:9000",
		S3Region:   "us-west-2",
		AccessKey:  "access",
		SecretKey:  "secret",
	})
	require.Equal(t, "http://localhost:9000", aws.ToString(opts.BaseEndpoint))
	require.True(t, opts.UsePathStyle)
	require.Equal(t, "us-west-2", opts.Region)
	creds, err := opts.Credentials.Retrieve(context.Background())
	require.NoError(t, err)
	require.Equal(t, "access", creds.AccessKeyID)
	require.Equal(t, "secret", creds.SecretAccessKey)
}

func TestParseLocationSplitsS3BucketAndKey(t *testing.T) {
	loc, err := E.UnwrapError(parseLocation("s3://bucket/runs/a"))
	require.NoError(t, err)
	require.Equal(t, schemeS3, locL.Scheme.Get(loc))
	require.Equal(t, "bucket", locL.Bucket.Get(loc))
	require.Equal(t, "runs/a", locL.Key.Get(loc))
}

func TestS3FileInputKeepsBucketRoot(t *testing.T) {
	loc, err := E.UnwrapError(parseLocation("s3://bucket/prefix/assemblies.zip"))
	require.NoError(t, err)
	got := fitInput(locL.Kind.Set(typeZip)(loc))
	require.Equal(t, "bucket", locL.Bucket.Get(got))
	require.Equal(t, "prefix/assemblies.zip", locL.Key.Get(got))
}

func TestDirPrefixExcludesSiblingPrefixes(t *testing.T) {
	require.Equal(t, "runs/a/", dirPrefix("runs/a"))
	require.Equal(t, "runs/a/", dirPrefix("runs/a/"))
	require.Empty(t, dirPrefix(""))
}
