package gbexport

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	En "github.com/IBM/fp-go/v2/endomorphism"
	ER "github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	PR "github.com/IBM/fp-go/v2/pair"
	P "github.com/IBM/fp-go/v2/predicate"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	endoutils "github.com/dictyBase/fp-go-loom/endomorphism"
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
	bundleL   = MakeFileBundleLenses()
	locL      = MakelocationLenses()
	sourceL   = MakesourceLenses()
	runL      = MakerunStateLenses()
	refL      = MakeobjectRefLenses()
	runParams = F.Pipe1(runL.State, L.Compose[runState](stateL.Params))
	runCtx    = F.Pipe1(runL.State, L.Compose[runState](stateL.Ctx))
	runRaw    = F.Pipe1(runParams, L.Compose[runState](paramsL.Input))
	runType   = F.Pipe1(runParams, L.Compose[runState](paramsL.InputType))
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
	// fileblobCreateDir reads and sets fileblob.Options.CreateDir.
	fileblobCreateDir = L.MakeLens(
		func(o fileblob.Options) bool { return o.CreateDir },
		func(o fileblob.Options, value bool) fileblob.Options { o.CreateDir = value; return o },
	)

	isDNAExt = predord.StrEq(dnaExt)
	isZipExt = predord.StrEq(zipExt)

	// locationOf builds a location from a parsed URL by its scheme.
	locationOf = MO.Match(
		E.Left[location](ErrUnsupportedScheme),
		[]O.Kleisli[*url.URL, E.Either[error, location]]{
			MO.Case(isLocalURL, localLocation),
			MO.Case(isS3URL, s3Location),
		},
	)
	// localLocationOf makes a file location with the URL path as bucket.
	localLocationOf = endoutils.BuildFrom([]func(*url.URL) En.Endomorphism[location]{
		F.Constant1[*url.URL](locL.Scheme.Set(schemeFile)),
		F.Flow2(urlPath, locL.Bucket.Set),
	})
	// s3LocationOf makes an s3 location with the URL host as bucket.
	s3LocationOf = endoutils.BuildFrom([]func(*url.URL) En.Endomorphism[location]{
		F.Constant1[*url.URL](locL.Scheme.Set(schemeS3)),
		F.Flow2(urlHost, locL.Bucket.Set),
		F.Flow3(urlPath, F.Bind2nd(strings.TrimPrefix, keySeparator), locL.Key.Set),
	})

	// kindOfExt maps a lower-case extension to an input kind.
	kindOfExt = MO.Match(typeDir, []O.Kleisli[string, string]{
		MO.Const(isDNAExt, typeDNA),
		MO.Const(isZipExt, typeZip),
	})
	isAutoInput     = F.Flow2(runType.Get, isAutoType)
	isExplicitInput = F.Flow2(runType.Get, isExplicitType)
	// inputKind infers the kind of an auto input and keeps an explicit one.
	inputKind = MO.Match(
		E.Left[string](ErrUnknownInputType),
		[]O.Kleisli[runState, E.Either[error, string]]{
			MO.Case(isAutoInput, autoInputKind),
			MO.Case(isExplicitInput, explicitInputKind),
		},
	)

	// fitInput roots a local file input at its parent directory.
	// Other inputs pass through unchanged.
	fitInput = F.Pipe1(isLocalFile, P.Fold(F.Identity[location], splitLocalFile))

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
	// s3Options sends the bucket name in the URL path so that MinIO
	// receives it. Only the static access and secret keys apply.
	s3Options = endoutils.BuildFrom([]func(Params) En.Endomorphism[s3.Options]{
		F.Flow2(paramsL.S3Endpoint.Get, s3Endpoint.Set),
		F.Flow2(paramsL.S3Region.Get, s3Region.Set),
		F.Flow2(staticCredentials, s3Credentials.Set),
		F.Constant1[Params](s3PathStyle.Set(true)),
	})

	targetL       = MakebucketTargetLenses()
	runInScheme   = F.Pipe1(runL.Input, L.Compose[runState](locL.Scheme))
	runInBucket   = F.Pipe1(runL.Input, L.Compose[runState](locL.Bucket))
	runOutScheme  = F.Pipe1(runL.Output, L.Compose[runState](locL.Scheme))
	runOutBucket  = F.Pipe1(runL.Output, L.Compose[runState](locL.Bucket))
	commonTargets = []func(runState) En.Endomorphism[bucketTarget]{
		F.Flow2(runCtx.Get, targetL.Ctx.Set),
		F.Flow4(runParams.Get, s3Options, newS3Client, targetL.Client.Set),
	}
	// inputTarget describes the input bucket of a run. It keeps the
	// zero fileblob.Options, which are the fileblob defaults.
	inputTarget = endoutils.BuildFrom(
		A.ArrayConcatAll(commonTargets, []func(runState) En.Endomorphism[bucketTarget]{
			F.Flow2(runInScheme.Get, targetL.Scheme.Set),
			F.Flow2(runInBucket.Get, targetL.Name.Set),
		}),
	)
	// outputTarget describes the output bucket of a run.
	outputTarget = endoutils.BuildFrom(
		A.ArrayConcatAll(commonTargets, []func(runState) En.Endomorphism[bucketTarget]{
			F.Flow2(runOutScheme.Get, targetL.Scheme.Set),
			F.Flow2(runOutBucket.Get, targetL.Name.Set),
			F.Constant1[runState](targetL.Options.Set(writeOptions())),
		}),
	)
	isLocalTarget = F.Flow2(targetL.Scheme.Get, predord.StrEq(schemeFile))
	isS3Target    = F.Flow2(targetL.Scheme.Get, predord.StrEq(schemeS3))
	// openTarget opens a directory bucket for a file target and an S3
	// bucket for an s3 target.
	openTarget = MO.Match(
		IOE.Left[*blob.Bucket](ErrUnsupportedScheme),
		[]O.Kleisli[bucketTarget, IOE.IOEither[error, *blob.Bucket]]{
			MO.Case(isLocalTarget, openLocalTarget),
			MO.Case(isS3Target, openS3Target),
		},
	)

	// inputRef points at the input key inside the open input bucket.
	inputRef = endoutils.BuildFrom([]func(runState) En.Endomorphism[objectRef]{
		F.Flow2(runL.InputBucket.Get, refL.Bucket.Set),
		F.Flow2(runCtx.Get, refL.Ctx.Set),
		F.Flow2(runKey.Get, refL.Key.Set),
	})
	// newSource makes a source from an object key.
	newSource = endoutils.BuildFrom([]func(string) En.Endomorphism[source]{
		sourceL.Key.Set,
		F.Flow2(stemOf, sourceL.Stem.Set),
	})

	isDNAInput = F.Flow2(runKind.Get, predord.StrEq(typeDNA))
	isZipInput = F.Flow2(runKind.Get, predord.StrEq(typeZip))
	isDirInput = F.Flow2(runKind.Get, predord.StrEq(typeDir))
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

// parseLocation parses raw as a local path or an s3:// URL.
func parseLocation(raw string) E.Either[error, location] {
	return F.Pipe3(
		raw,
		E.Eitherize1(url.Parse),
		E.MapLeft[*url.URL](ER.OnError("parse location")),
		E.Chain(locationOf),
	)
}

// localLocation returns the file location of u.
func localLocation(u *url.URL) E.Either[error, location] {
	return F.Pipe2(u, localLocationOf, E.Right[error, location])
}

// s3Location returns the s3 location of u.
func s3Location(u *url.URL) E.Either[error, location] {
	return F.Pipe2(u, s3LocationOf, E.Right[error, location])
}

// autoKind infers the input kind from the extension of raw.
func autoKind(raw string) string {
	return F.Pipe4(
		raw,
		F.Bind2nd(strings.TrimSuffix, keySeparator),
		path.Ext,
		strings.ToLower,
		kindOfExt,
	)
}

// autoInputKind infers the kind from the raw input location of r.
func autoInputKind(r runState) E.Either[error, string] {
	return F.Pipe3(r, runRaw.Get, autoKind, E.Right[error, string])
}

// explicitInputKind returns the input type flag of r.
func explicitInputKind(r runState) E.Either[error, string] {
	return F.Pipe2(r, runType.Get, E.Right[error, string])
}

// splitLocalFile turns the local file path held in Bucket into its
// parent directory (Bucket) and its file name (Key). The path is read
// once, so the order of the two Set steps does not matter.
func splitLocalFile(loc location) location {
	file := locL.Bucket.Get(loc)
	return F.Pipe2(
		loc,
		locL.Bucket.Set(filepath.Dir(file)),
		locL.Key.Set(filepath.Base(file)),
	)
}

// resolveLocations parses the raw input and output locations of r,
// infers the input kind, and roots a local file input at its parent
// directory. It returns r with Input and Output set, or the first
// parse or input-type error.
func resolveLocations(r runState) E.Either[error, runState] {
	return F.Pipe5(
		r,
		E.Of[error, runState],
		E.Bind(runL.Input.Set, F.Flow2(runRaw.Get, parseLocation)),
		E.Bind(runKind.Set, inputKind),
		E.Map[error](runL.Input.Modify(fitInput)),
		E.Bind(runL.Output.Set, F.Flow2(runOutRaw.Get, parseLocation)),
	)
}

// openInput opens the input bucket of r.
func openInput(r runState) IOE.IOEither[error, *blob.Bucket] {
	return F.Pipe2(r, inputTarget, openTarget)
}

// openOutput opens the output bucket of r.
func openOutput(r runState) IOE.IOEither[error, *blob.Bucket] {
	return F.Pipe2(r, outputTarget, openTarget)
}

// singleSource returns the one .dna file input as the only source.
func singleSource(r runState) IOE.IOEither[error, []source] {
	return F.Pipe4(r, runKey.Get, newSource, A.Of[source], IOE.Of[error, []source])
}

// archiveSource reads the zip input and returns its selectable members.
func archiveSource(r runState) IOE.IOEither[error, []source] {
	return F.Pipe4(
		r,
		inputRef,
		readRef,
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
		inputRef,
		refL.Key.Modify(dirPrefix),
		listTwoLevels,
		IOE.Map[error](objectSources),
	)
}

// objectSources makes a source for every selectable listed object.
func objectSources(objects []*blob.ListObject) []source {
	return F.Pipe2(objects, A.Filter(selectableObject), A.Map(objectSource))
}

// objectSource makes a source from the key of a listed object.
func objectSource(object *blob.ListObject) source {
	return F.Pipe2(object, objectKey, newSource)
}

// dirPrefix ends a non-empty directory key with one separator.
func dirPrefix(key string) string {
	return F.Pipe2(
		key,
		F.Bind2nd(strings.TrimSuffix, keySeparator),
		strutils.MapNonEmpty(S.Append(keySeparator)),
	)
}

// reserveStems sorts the sources by key and flags every source whose
// stem appears earlier in that order.
func reserveStems(sources []source) []source {
	return F.Pipe4(
		sources,
		A.SortByKey(S.Ord, sourceL.Key.Get),
		PR.Of[[]source],
		PR.MapTail[[]source](A.MapWithIndex(PR.MakePair[int, source])),
		PR.Paired(markDuplicates),
	)
}

// writeOptions configures the output bucket. It makes the output
// directory when it does not exist, and it writes only the artifact
// files, with no ".attrs" metadata file next to each one.
func writeOptions() fileblob.Options {
	return F.Pipe2(
		fileblob.Options{},
		fileblobCreateDir.Set(true),
		skipMetadataFiles,
	)
}

// skipMetadataFiles stops fileblob from writing a ".attrs" metadata
// file next to each object. It is a function and not a lens because
// the type of fileblob.Options.Metadata is unexported.
func skipMetadataFiles(o fileblob.Options) fileblob.Options {
	o.Metadata = fileblob.MetadataDontWrite
	return o
}

func isLocalURL(u *url.URL) bool { return u.Scheme == "" || u.Scheme == schemeFile }

func isAutoType(kind string) bool { return kind == "" || kind == typeAuto }

func isExplicitType(kind string) bool {
	return kind == typeDNA || kind == typeZip || kind == typeDir
}

func isLocalFile(loc location) bool {
	return locL.Scheme.Get(loc) == schemeFile && locL.Kind.Get(loc) != typeDir
}

func isS3URL(u *url.URL) bool { return u.Scheme == schemeS3 }

func urlHost(u *url.URL) string { return u.Host }

func urlPath(u *url.URL) string { return u.Path }

func staticCredentials(params Params) aws.CredentialsProvider {
	key, secret := paramsL.AccessKey.Get(params), paramsL.SecretKey.Get(params)
	return credentials.NewStaticCredentialsProvider(key, secret, "")
}

// newS3Client builds the client before any raw bucket effect.
//
// fp-go:allow-trivial-delegate s3.New is variadic and cannot be a unary pipeline step
func newS3Client(opts s3.Options) *s3.Client { return s3.New(opts) }

func openLocalTarget(target bucketTarget) IOE.IOEither[error, *blob.Bucket] {
	return F.Pipe1(
		IOE.TryCatchError(
			func() (*blob.Bucket, error) {
				dir, opts := targetL.Name.Get(target), targetL.Options.Get(target)
				return fileblob.OpenBucket(dir, &opts)
			},
		),
		IOE.MapLeft[*blob.Bucket](ER.OnError("open bucket %q", targetL.Name.Get(target))),
	)
}

// openS3Target never adds the credentials to an error or a log record.
func openS3Target(target bucketTarget) IOE.IOEither[error, *blob.Bucket] {
	return F.Pipe1(
		IOE.TryCatchError(
			func() (*blob.Bucket, error) {
				ctx, client := targetL.Ctx.Get(target), targetL.Client.Get(target)
				name := targetL.Name.Get(target)
				return s3blob.OpenBucket(ctx, client, name, nil)
			},
		),
		IOE.MapLeft[*blob.Bucket](ER.OnError("open S3 bucket %q", targetL.Name.Get(target))),
	)
}

func closeResource[R io.Closer](resource R) IOE.IOEither[error, struct{}] {
	return IOE.TryCatchError(
		func() (struct{}, error) { return struct{}{}, resource.Close() },
	)
}

// withResource keeps an opened resource alive for use and closes it on
// success and on failure.
func withResource[R io.Closer, A any](
	open IOE.IOEither[error, R],
	use func(R) IOE.IOEither[error, A],
) IOE.IOEither[error, A] {
	return F.Pipe1(use, IOE.WithResource[A](open, closeResource[R]))
}

func readRef(ref objectRef) IOE.IOEither[error, []byte] {
	return F.Pipe1(
		IOE.TryCatchError(
			func() ([]byte, error) {
				bucket, ctx, key := refL.Bucket.Get(ref), refL.Ctx.Get(ref), refL.Key.Get(ref)
				return bucket.ReadAll(ctx, key)
			},
		),
		IOE.MapLeft[[]byte](ER.OnError("read %q", refL.Key.Get(ref))),
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

// selectableMember keeps .dna files at the archive root and one folder deep.
func selectableMember(member *zip.File) bool {
	name := memberName(member)
	depth, base := strings.Count(name, keySeparator), path.Base(name)
	return !member.FileInfo().IsDir() && depth <= maxFolderDepth &&
		!strings.HasPrefix(name, macMetadataDir+keySeparator) && selectableName(base)
}

func selectableName(base string) bool {
	ext := path.Ext(base)
	return !strings.HasPrefix(base, appleDoublePrefix) && strings.EqualFold(ext, dnaExt)
}

func stemOf(key string) string {
	base := path.Base(key)
	return strings.TrimSuffix(base, path.Ext(base))
}

// listTwoLevels lists the prefix and each first-level folder once.
func listTwoLevels(ref objectRef) IOE.IOEither[error, []*blob.ListObject] {
	return F.Pipe3(
		ref,
		listObjects,
		IOE.Chain(F.Bind1st(expandFolders, ref)),
		IOE.MapLeft[[]*blob.ListObject](ER.OnError("list input %q", refL.Key.Get(ref))),
	)
}

func expandFolders(
	ref objectRef,
	objects []*blob.ListObject,
) IOE.IOEither[error, []*blob.ListObject] {
	return F.Pipe3(
		objects,
		A.Filter(isListedFolder),
		IOE.TraverseArraySeq(F.Flow2(F.Bind1st(childRef, ref), listObjects)),
		IOE.Map[error](F.Flow2(A.Flatten[*blob.ListObject], A.Concat(objects))),
	)
}

func childRef(ref objectRef, folder *blob.ListObject) objectRef {
	return F.Pipe3(folder, objectKey, refL.Key.Set, En.Read(ref))
}

// listObjects drains one delimited listing of the reference prefix.
//
// fp-go:allow-trycatch-iferr blob iterator drain separates io.EOF from failures in one raw effect
func listObjects(ref objectRef) IOE.IOEither[error, []*blob.ListObject] {
	return IOE.TryCatchError(
		func() ([]*blob.ListObject, error) {
			bucket, ctx, prefix := refL.Bucket.Get(ref), refL.Ctx.Get(ref), refL.Key.Get(ref)
			opts := listOptions(prefix)
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

func isListedFolder(object *blob.ListObject) bool {
	key := objectKey(object)
	folder := strings.TrimSuffix(key, keySeparator)
	return object.IsDir && path.Base(folder) != macMetadataDir
}

func selectableObject(object *blob.ListObject) bool {
	key := objectKey(object)
	base := path.Base(key)
	return !object.IsDir && selectableName(base)
}

// markDuplicates flags each indexed source whose stem appears before its
// index in sorted.
func markDuplicates(sorted []source, indexed []PR.Pair[int, source]) []source {
	return F.Pipe1(indexed, A.Map(F.Bind1st(markDuplicate, sorted)))
}

// markDuplicate flags the source of item when a source before its index
// in sorted has the same stem.
func markDuplicate(sorted []source, item PR.Pair[int, source]) source {
	return F.Pipe3(
		sorted[:PR.Head(item)],
		A.Any(F.Bind1st(sameStem, PR.Tail(item))),
		sourceL.Duplicate.Set,
		En.Read(PR.Tail(item)),
	)
}

func sameStem(first, second source) bool {
	return sourceL.Stem.Get(first) == sourceL.Stem.Get(second)
}
