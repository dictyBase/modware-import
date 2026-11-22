package gbexport

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"io"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	ER "github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	ORD "github.com/IBM/fp-go/v2/ord"
	P "github.com/IBM/fp-go/v2/predicate"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	MO "github.com/dictyBase/fp-go-loom/matchopt"
	predarrays "github.com/dictyBase/fp-go-loom/predicate/array"
	predord "github.com/dictyBase/fp-go-loom/predicate/ord"
	"github.com/dictyBase/fp-go-loom/strutils"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
	"gocloud.dev/blob/s3blob"
)

var (
	paramsL   = MakeParamsLenses()
	stateL    = MakeStateLenses()
	locL      = MakelocationLenses()
	sourceL   = MakesourceLenses()
	runL      = MakerunStateLenses()
	runParams = F.Pipe1(runL.State, L.Compose[runState](stateL.Params))
	runCtx    = F.Pipe1(runL.State, L.Compose[runState](stateL.Ctx))
	runRaw    = F.Pipe1(runParams, L.Compose[runState](paramsL.Input))
	runOutRaw = F.Pipe1(runParams, L.Compose[runState](paramsL.Output))
	runKind   = F.Pipe1(runL.Input, L.Compose[runState](locL.Kind))
	runKey    = F.Pipe1(runL.Input, L.Compose[runState](locL.Key))

	listPrefix = L.MakeLens(
		func(o blob.ListOptions) string { return o.Prefix },
		func(o blob.ListOptions, value string) blob.ListOptions { o.Prefix = value; return o },
	)
	listDelimiter = L.MakeLens(
		func(o blob.ListOptions) string { return o.Delimiter },
		func(o blob.ListOptions, value string) blob.ListOptions { o.Delimiter = value; return o },
	)
	fileblobCreateDir = L.MakeLens(
		func(o fileblob.Options) bool { return o.CreateDir },
		func(o fileblob.Options, value bool) fileblob.Options { o.CreateDir = value; return o },
	)
	s3Endpoint = L.MakeLens(
		func(o s3.Options) string { return aws.ToString(o.BaseEndpoint) },
		func(o s3.Options, value string) s3.Options { o.BaseEndpoint = aws.String(value); return o },
	)
	s3Region = L.MakeLens(
		func(o s3.Options) string { return o.Region },
		func(o s3.Options, value string) s3.Options { o.Region = value; return o },
	)
	s3Credentials = L.MakeLens(
		func(o s3.Options) aws.CredentialsProvider { return o.Credentials },
		func(o s3.Options, value aws.CredentialsProvider) s3.Options { o.Credentials = value; return o },
	)
	s3PathStyle = L.MakeLens(
		func(o s3.Options) bool { return o.UsePathStyle },
		func(o s3.Options, value bool) s3.Options { o.UsePathStyle = value; return o },
	)

	isDNAExt     = predord.StrEq(dnaExt)
	isZipExt     = predord.StrEq(zipExt)
	isDNAKind    = predord.StrEq(typeDNA)
	isZipKind    = predord.StrEq(typeZip)
	isDirKind    = predord.StrEq(typeDir)
	isFileScheme = predord.StrEq(schemeFile)
	isS3Scheme   = predord.StrEq(schemeS3)
	intLeq       = ORD.Leq(predord.IntOrd)

	hasS3Prefix    = F.Bind2nd(strings.HasPrefix, schemeS3+schemeSeparator)
	hasFilePrefix  = F.Bind2nd(strings.HasPrefix, schemeFile+schemeSeparator)
	hasOtherScheme = F.Bind2nd(strings.Contains, schemeSeparator)

	isFileLocation = F.Flow2(locL.Scheme.Get, isFileScheme)
	isDirLocation  = F.Flow2(locL.Kind.Get, isDirKind)
	// isLocalFile accepts a local input that is not a directory.
	isLocalFile = F.Pipe1(isFileLocation, P.And(P.Not(isDirLocation)))

	// isAppleDouble matches the "._" resource files that macOS adds.
	isAppleDouble = F.Bind2nd(strings.HasPrefix, appleDoublePrefix)
	hasDNAExt     = F.Flow2(path.Ext, F.Bind2nd(strings.EqualFold, dnaExt))
	// selectableName accepts a .dna file name that is not an AppleDouble file.
	selectableName = F.Pipe1(hasDNAExt, P.And(P.Not(isAppleDouble)))

	withinFolderDepth = F.Flow3(
		memberName,
		F.Bind2nd(strings.Count, keySeparator),
		intLeq(maxFolderDepth),
	)
	inMacMetadata = F.Flow2(
		memberName,
		F.Bind2nd(strings.HasPrefix, macMetadataDir+keySeparator),
	)
	hasSelectableBase = F.Flow3(memberName, path.Base, selectableName)
	// selectableMember keeps .dna files at the archive root and one folder deep.
	selectableMember = F.Pipe3(
		P.Not(isMemberDir),
		P.And(withinFolderDepth),
		P.And(P.Not(inMacMetadata)),
		P.And(hasSelectableBase),
	)

	isMacMetadataName   = predord.StrEq(macMetadataDir)
	isMacMetadataFolder = F.Flow4(
		objectKey,
		F.Bind2nd(strings.TrimSuffix, keySeparator),
		path.Base,
		isMacMetadataName,
	)
	// isListedFolder accepts a listed folder that is not __MACOSX.
	isListedFolder   = F.Pipe1(isObjectDir, P.And(P.Not(isMacMetadataFolder)))
	hasSelectableKey = F.Flow3(objectKey, path.Base, selectableName)
	// selectableObject accepts a listed .dna file.
	selectableObject = F.Pipe1(P.Not(isObjectDir), P.And(hasSelectableKey))

	// parseLocation classifies raw by URL scheme prefix before parsing URLs,
	// so literal local paths with '#', '?', or '%' pass through unchanged.
	parseLocation = MO.Match(
		E.Left[location](ErrUnsupportedScheme),
		[]O.Kleisli[string, E.Either[error, location]]{
			MO.Case(hasS3Prefix, parseS3Location),
			MO.Case(hasFilePrefix, parseFileURLLocation),
			MO.Const(hasOtherScheme, E.Left[location](ErrUnsupportedScheme)),
			MO.Default(localPathLocation),
		},
	)

	// kindOfExt maps a lower-case extension to an input kind.
	kindOfExt = MO.Match(typeDir, []O.Kleisli[string, string]{
		MO.Const(isDNAExt, typeDNA),
		MO.Const(isZipExt, typeZip),
	})

	// fitInput roots a local file input at its parent directory.
	// Other inputs pass through unchanged.
	fitInput = F.Pipe1(isLocalFile, P.Fold(F.Identity[location], splitLocalFile))

	isDNAInput = F.Flow2(runKind.Get, isDNAKind)
	isZipInput = F.Flow2(runKind.Get, isZipKind)
	isDirInput = F.Flow2(runKind.Get, isDirKind)
	// selectSources lists the sources of a single file, an archive, or a
	// directory input.
	selectSources = MO.Match(
		IOE.Left[[]source](ErrUnknownInputType),
		[]O.Kleisli[runState, IOE.IOEither[error, []source]]{
			MO.Case(isDNAInput, singleSource),
			MO.Case(isZipInput, archiveSource),
			MO.Case(isDirInput, dirSource),
		},
	)
	requireSources = E.FromPredicate(
		F.Flow2(runL.Sources.Get, predarrays.IsNonEmpty[source]()),
		F.Constant1[runState](ErrNoDNAFiles),
	)
)

// localPathLocation wraps a literal filesystem path as a file location.
func localPathLocation(raw string) E.Either[error, location] {
	return F.Pipe3(
		location{},
		locL.Scheme.Set(schemeFile),
		locL.Bucket.Set(raw),
		E.Right[error, location],
	)
}

func parseURL(raw string) E.Either[error, *url.URL] {
	return F.Pipe2(
		raw,
		E.Eitherize1(url.Parse),
		E.MapLeft[*url.URL](ER.OnError("parse location")),
	)
}

// parseFileURLLocation parses a file:// URL into a file location.
func parseFileURLLocation(raw string) E.Either[error, location] {
	return F.Pipe2(raw, parseURL, E.Map[error](fileURLLocation))
}

func fileURLLocation(u *url.URL) location {
	return F.Pipe2(
		location{},
		locL.Scheme.Set(schemeFile),
		locL.Bucket.Set(u.Path),
	)
}

// parseS3Location parses an s3:// URL into an S3 location.
func parseS3Location(raw string) E.Either[error, location] {
	return F.Pipe2(raw, parseURL, E.Map[error](s3URLLocation))
}

func s3URLLocation(u *url.URL) location {
	key := strings.TrimPrefix(u.Path, keySeparator)
	return F.Pipe3(
		location{},
		locL.Scheme.Set(schemeS3),
		locL.Bucket.Set(u.Host),
		locL.Key.Set(key),
	)
}

// inputKind infers the input kind from the extension of raw: .dna is
// a single file, .zip is an archive, and anything else is a folder.
func inputKind(raw string) string {
	return F.Pipe4(
		raw,
		F.Bind2nd(strings.TrimSuffix, keySeparator),
		path.Ext,
		strings.ToLower,
		kindOfExt,
	)
}

// splitLocalFile turns the local file path held in Bucket into its
// parent directory (Bucket) and its file name (Key).
func splitLocalFile(loc location) location {
	file := locL.Bucket.Get(loc)
	dir, base := filepath.Dir(file), filepath.Base(file)
	return F.Pipe2(
		loc,
		locL.Bucket.Set(dir),
		locL.Key.Set(base),
	)
}

// resolveLocations parses the raw input and output locations of r,
// infers the input kind, and roots a local file input at its parent
// directory.
func resolveLocations(r runState) E.Either[error, runState] {
	return F.Pipe5(
		r,
		E.Of[error, runState],
		E.Bind(runL.Input.Set, F.Flow2(runRaw.Get, parseLocation)),
		E.Let[error](runKind.Set, F.Flow2(runRaw.Get, inputKind)),
		E.Map[error](runL.Input.Modify(fitInput)),
		E.Bind(runL.Output.Set, F.Flow2(runOutRaw.Get, parseLocation)),
	)
}

// s3Options sends the bucket name in the URL path so MinIO receives it.
func s3Options(params Params) s3.Options {
	endpoint := paramsL.S3Endpoint.Get(params)
	region := paramsL.S3Region.Get(params)
	creds := staticCredentials(params)
	return F.Pipe4(
		s3.Options{},
		s3Endpoint.Set(endpoint),
		s3Region.Set(region),
		s3Credentials.Set(creds),
		s3PathStyle.Set(true),
	)
}

// openInputBucket opens the input bucket with default fileblob options.
func openInputBucket(r runState) IOE.IOEither[error, *blob.Bucket] {
	loc := runL.Input.Get(r)
	return openBucket(r, loc, fileblob.Options{})
}

// openOutputBucket opens the output bucket with directory creation enabled.
func openOutputBucket(r runState) IOE.IOEither[error, *blob.Bucket] {
	loc := runL.Output.Get(r)
	opts := writeOptions()
	return openBucket(r, loc, opts)
}

func openBucket(
	r runState,
	loc location,
	opts fileblob.Options,
) IOE.IOEither[error, *blob.Bucket] {
	scheme := locL.Scheme.Get(loc)
	name := locL.Bucket.Get(loc)
	ctx, params := runCtx.Get(r), runParams.Get(r)
	localIOE := openLocalBucket(name, opts)
	s3IOE := openS3Bucket(ctx, params, name)
	return F.Pipe1(
		scheme,
		MO.Match(
			IOE.Left[*blob.Bucket](ErrUnsupportedScheme),
			[]O.Kleisli[string, IOE.IOEither[error, *blob.Bucket]]{
				MO.Const(isFileScheme, localIOE),
				MO.Const(isS3Scheme, s3IOE),
			},
		),
	)
}

func openLocalBucket(dir string, opts fileblob.Options) IOE.IOEither[error, *blob.Bucket] {
	return F.Pipe1(
		IOE.TryCatchError(
			func() (*blob.Bucket, error) {
				return fileblob.OpenBucket(dir, &opts)
			},
		),
		IOE.MapLeft[*blob.Bucket](ER.OnError("open bucket %q", dir)),
	)
}

// openS3Bucket builds the S3 client lazily so local runs never allocate one
// and credentials never appear in errors or logs.
func openS3Bucket(
	ctx context.Context,
	params Params,
	name string,
) IOE.IOEither[error, *blob.Bucket] {
	return F.Pipe1(
		IOE.TryCatchError(
			func() (*blob.Bucket, error) {
				opts := s3Options(params)
				client := s3.New(opts)
				return s3blob.OpenBucket(ctx, client, name, nil)
			},
		),
		IOE.MapLeft[*blob.Bucket](ER.OnError("open S3 bucket %q", name)),
	)
}

// singleSource returns the one .dna file input as the only source.
func singleSource(r runState) IOE.IOEither[error, []source] {
	return F.Pipe4(r, runKey.Get, newSource, A.Of[source], IOE.Of[error, []source])
}

// archiveSource reads the zip input and returns its selectable members.
func archiveSource(r runState) IOE.IOEither[error, []source] {
	return F.Pipe4(
		r,
		runKey.Get,
		F.Bind1st(readBucketObject, r),
		IOE.ChainEitherK(openArchive),
		IOE.Map[error](archiveSources),
	)
}

// archiveSources makes a source for every selectable archive member.
func archiveSources(archive *zip.Reader) []source {
	return F.Pipe3(archive, archiveFiles, A.Filter(selectableMember), A.Map(memberSource))
}

// memberSource makes a source for one archive member.
func memberSource(member *zip.File) source {
	return F.Pipe3(member, memberName, newSource, sourceL.Member.Set(member))
}

// dirSource lists the directory input two levels deep and returns its
// selectable objects.
func dirSource(r runState) IOE.IOEither[error, []source] {
	return F.Pipe4(
		r,
		runKey.Get,
		dirPrefix,
		F.Bind1st(listTwoLevels, r),
		IOE.Map[error](objectSources),
	)
}

// objectSources makes a source for every selectable listed object.
func objectSources(objects []*blob.ListObject) []source {
	return F.Pipe2(objects, A.Filter(selectableObject), A.Map(F.Flow2(objectKey, newSource)))
}

func newSource(key string) source {
	stem := stemOf(key)
	return F.Pipe2(
		source{},
		sourceL.Key.Set(key),
		sourceL.Stem.Set(stem),
	)
}

// dirPrefix ends a non-empty directory key with one separator.
func dirPrefix(key string) string {
	return F.Pipe2(
		key,
		F.Bind2nd(strings.TrimSuffix, keySeparator),
		strutils.MapNonEmpty(S.Append(keySeparator)),
	)
}

// writeOptions configures the output bucket to create missing directories
// and omit sidecar ".attrs" metadata files.
func writeOptions() fileblob.Options {
	return F.Pipe2(
		fileblob.Options{},
		fileblobCreateDir.Set(true),
		skipMetadataFiles,
	)
}

// skipMetadataFiles disables ".attrs" sidecar files. It is a function
// rather than a lens because fileblob.Options.Metadata has an unexported type.
func skipMetadataFiles(o fileblob.Options) fileblob.Options {
	o.Metadata = fileblob.MetadataDontWrite
	return o
}

func staticCredentials(params Params) aws.CredentialsProvider {
	key, secret := paramsL.AccessKey.Get(params), paramsL.SecretKey.Get(params)
	return credentials.NewStaticCredentialsProvider(key, secret, "")
}

// closeResource closes an opened bucket or archive member.
func closeResource[R io.Closer](resource R) IOE.IOEither[error, struct{}] {
	return IOE.TryCatchError(
		func() (struct{}, error) { return struct{}{}, resource.Close() },
	)
}

func readBucketObject(r runState, key string) IOE.IOEither[error, []byte] {
	bucket, ctx := runL.InputBucket.Get(r), runCtx.Get(r)
	return F.Pipe1(
		IOE.TryCatchError(
			func() ([]byte, error) {
				return bucket.ReadAll(ctx, key)
			},
		),
		IOE.MapLeft[[]byte](ER.OnError("read %q", key)),
	)
}

// openArchive opens data as a zip archive.
func openArchive(data []byte) E.Either[error, *zip.Reader] {
	return F.Pipe2(data, newZipReader, E.MapLeft[*zip.Reader](ER.OnError("open archive")))
}

func newZipReader(data []byte) E.Either[error, *zip.Reader] {
	reader, size := bytes.NewReader(data), int64(len(data))
	archive, err := zip.NewReader(reader, size)
	return E.TryCatchError(archive, err)
}

func archiveFiles(archive *zip.Reader) []*zip.File { return archive.File }

func memberName(member *zip.File) string { return member.Name }

func isMemberDir(member *zip.File) bool { return member.FileInfo().IsDir() }

func stemOf(key string) string {
	base := path.Base(key)
	ext := path.Ext(base)
	return strings.TrimSuffix(base, ext)
}

// listTwoLevels lists the prefix and each first-level folder once,
// stopping at the first error.
func listTwoLevels(r runState, prefix string) IOE.IOEither[error, []*blob.ListObject] {
	return F.Pipe3(
		prefix,
		F.Bind1st(listObjects, r),
		IOE.Chain(F.Bind1st(expandFolders, r)),
		IOE.MapLeft[[]*blob.ListObject](ER.OnError("list input %q", prefix)),
	)
}

// expandFolders lists each first-level subfolder sequentially and stops on
// the first failure.
func expandFolders(
	r runState,
	objects []*blob.ListObject,
) IOE.IOEither[error, []*blob.ListObject] {
	return F.Pipe3(
		objects,
		A.Filter(isListedFolder),
		A.Map(objectKey),
		A.Reduce(r.appendFolderObjects, IOE.Of[error](objects)),
	)
}

func (r runState) appendFolderObjects(
	acc IOE.IOEither[error, []*blob.ListObject],
	prefix string,
) IOE.IOEither[error, []*blob.ListObject] {
	return F.Pipe1(acc, IOE.Chain(F.Bind2nd(r.listAndAppend, prefix)))
}

func (r runState) listAndAppend(
	existing []*blob.ListObject,
	prefix string,
) IOE.IOEither[error, []*blob.ListObject] {
	return F.Pipe2(
		prefix,
		F.Bind1st(listObjects, r),
		IOE.Map[error](A.Concat(existing)),
	)
}

// listObjects drains one delimited listing of prefix.
//
// fp-go:allow-trycatch-iferr blob iterator drain separates io.EOF from failures in one raw effect
func listObjects(r runState, prefix string) IOE.IOEither[error, []*blob.ListObject] {
	bucket, ctx := runL.InputBucket.Get(r), runCtx.Get(r)
	opts := listOptions(prefix)
	return IOE.TryCatchError(
		func() ([]*blob.ListObject, error) {
			iter := bucket.List(&opts)
			var objects []*blob.ListObject
			for {
				object, err := iter.Next(ctx)
				if errors.Is(err, io.EOF) {
					return objects, nil
				}
				if err != nil {
					return nil, err
				}
				objects = append(objects, object)
			}
		},
	)
}

// listOptions lists one level under prefix: the delimiter groups
// deeper keys into folder entries.
func listOptions(prefix string) blob.ListOptions {
	return F.Pipe2(
		blob.ListOptions{},
		listDelimiter.Set(keySeparator),
		listPrefix.Set(prefix),
	)
}

func objectKey(object *blob.ListObject) string { return object.Key }

func isObjectDir(object *blob.ListObject) bool { return object.IsDir }
