package gbexport

import (
	"context"
	"errors"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/require"
	"gocloud.dev/blob"
	"gocloud.dev/blob/driver"
	"gocloud.dev/gcerrors"
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

func TestParseLocationHandlesSpecialLocalCharsAndFileURL(t *testing.T) {
	for _, raw := range []string{
		"/tmp/construct#1.dna",
		"/tmp/dir?v=2/a.dna",
		"/tmp/100%valid/a.dna",
	} {
		loc, err := E.UnwrapError(parseLocation(raw))
		require.NoError(t, err, raw)
		require.Equal(t, schemeFile, locL.Scheme.Get(loc))
		require.Equal(t, raw, locL.Bucket.Get(loc))
	}

	loc, err := E.UnwrapError(parseLocation("file:///tmp/my%20dir/a.dna"))
	require.NoError(t, err)
	require.Equal(t, schemeFile, locL.Scheme.Get(loc))
	require.Equal(t, "/tmp/my dir/a.dna", locL.Bucket.Get(loc))

	_, err = E.UnwrapError(parseLocation("https://example.org/a.dna"))
	require.ErrorIs(t, err, ErrUnsupportedScheme)
}

type failingListDriver struct {
	driver.Bucket
	prefixes []string
	failOn   string
	err      error
}

func (d *failingListDriver) ErrorCode(err error) gcerrors.ErrorCode {
	if errors.Is(err, d.err) {
		return gcerrors.Internal
	}
	return gcerrors.Unknown
}

func (d *failingListDriver) ListPaged(
	_ context.Context,
	opts *driver.ListOptions,
) (*driver.ListPage, error) {
	d.prefixes = append(d.prefixes, opts.Prefix)
	if opts.Prefix == d.failOn {
		return nil, d.err
	}
	if opts.Prefix == "" {
		return &driver.ListPage{
			Objects: []*driver.ListObject{
				{Key: "a/", IsDir: true},
				{Key: "b/", IsDir: true},
			},
		}, nil
	}
	return &driver.ListPage{}, nil
}

func (d *failingListDriver) Close() error { return nil }

func TestListTwoLevelsStopsOnFirstFolderError(t *testing.T) {
	listErr := errors.New("subfolder list failed")
	drv := &failingListDriver{failOn: "a/", err: listErr}
	bucket := blob.NewBucket(drv)
	t.Cleanup(func() { _ = bucket.Close() })

	state := runState{
		State:       State{Ctx: context.Background()},
		InputBucket: bucket,
	}
	_, err := E.UnwrapError(listTwoLevels(state, "")())
	require.ErrorIs(t, err, listErr)
	require.Equal(t, []string{"", "a/"}, drv.prefixes)
}
