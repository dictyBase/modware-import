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
	M "github.com/IBM/fp-go/v2/monoid"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	PR "github.com/IBM/fp-go/v2/pair"
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
	createDir = L.MakeLens(
		func(o fileblob.Options) bool { return o.CreateDir },
		func(o fileblob.Options, value bool) fileblob.Options { o.CreateDir = value; return o },
	)

	// readOptions leaves the local input bucket at the fileblob defaults.
	readOptions *fileblob.Options
	// writeOptions creates the output root and writes no metadata sidecars.
	writeOptions = new(F.Pipe2(fileblob.Options{}, createDir.Set(true), skipSidecars))

	parseLocation = F.Flow3(
		E.Eitherize1(url.Parse),
		E.MapLeft[*url.URL](ER.OnError("parse location")),
		E.Chain(MO.Match(
			E.Left[location](ErrUnsupportedScheme),
			[]O.Kleisli[*url.URL, E.Either[error, location]]{
				MO.Case(isLocalURL, F.Flow2(localLocation, E.Right[error, location])),
				MO.Case(isS3URL, F.Flow2(s3Location, E.Right[error, location])),
			},
		)),
	)
	localLocation = F.Bind2nd(build[*url.URL, location], []func(*url.URL) En.Endomorphism[location]{
		F.Constant1[*url.URL](locL.Scheme.Set(schemeFile)),
		F.Flow2(urlPath, locL.Bucket.Set),
	})
	s3Location = F.Bind2nd(build[*url.URL, location], []func(*url.URL) En.Endomorphism[location]{
		F.Constant1[*url.URL](locL.Scheme.Set(schemeS3)),
		F.Flow2(urlHost, locL.Bucket.Set),
		F.Flow3(urlPath, F.Bind2nd(strings.TrimPrefix, keySeparator), locL.Key.Set),
	})
	isDNAExt = predord.StrEq(dnaExt)
	isZipExt = predord.StrEq(zipExt)
	autoKind = F.Flow4(
		F.Bind2nd(strings.TrimSuffix, keySeparator),
		path.Ext,
		strings.ToLower,
		MO.Match(typeDir, []O.Kleisli[string, string]{
			MO.Const(isDNAExt, typeDNA),
			MO.Const(isZipExt, typeZip),
		}),
	)
	inputKind = F.Flow3(
		PR.Of[runState],
		PR.BiMap(runType.Get, runRaw.Get),
		MO.Match(
			E.Left[string](ErrUnknownInputType),
			[]O.Kleisli[PR.Pair[string, string], E.Either[error, string]]{
				MO.Case(
					F.Flow2(PR.Head[string, string], isAutoType),
					F.Flow3(PR.Tail[string, string], autoKind, E.Right[error, string]),
				),
				MO.Case(
					F.Flow2(PR.Head[string, string], isExplicitType),
					F.Flow2(PR.Head[string, string], E.Right[error, string]),
				),
			},
		),
	)
	// fitInput roots a local file input at its parent directory.
	fitInput = F.Pipe1(
		isLocalFile,
		P.Fold(
			F.Identity[location],
			F.Bind2nd(update[location], []func(location) En.Endomorphism[location]{
				F.Flow3(locL.Bucket.Get, filepath.Dir, locL.Bucket.Set),
				F.Flow3(locL.Bucket.Get, filepath.Base, locL.Key.Set),
			}),
		),
	)
	resolveLocations = F.Flow5(
		E.Of[error, runState],
		E.Bind(runL.Input.Set, F.Flow2(runRaw.Get, parseLocation)),
		E.Bind(runKind.Set, inputKind),
		E.Map[error](runL.Input.Modify(fitInput)),
		E.Bind(runL.Output.Set, F.Flow2(runOutRaw.Get, parseLocation)),
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
	// s3Options sends the bucket name in the URL path so that MinIO
	// receives it. Only the static access and secret keys apply.
	s3Options = F.Bind2nd(build[Params, s3.Options], []func(Params) En.Endomorphism[s3.Options]{
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
	openTarget = MO.Match(
		IOE.Left[*blob.Bucket](ErrUnsupportedScheme),
		[]O.Kleisli[bucketTarget, IOE.IOEither[error, *blob.Bucket]]{
			MO.Case(F.Flow2(targetL.Scheme.Get, predord.StrEq(schemeFile)), openLocalTarget),
			MO.Case(F.Flow2(targetL.Scheme.Get, predord.StrEq(schemeS3)), openS3Target),
		},
	)
	openInput = F.Flow2(
		F.Bind2nd(
			build[runState, bucketTarget],
			A.ArrayConcatAll(commonTargets, []func(runState) En.Endomorphism[bucketTarget]{
				F.Flow2(runInScheme.Get, targetL.Scheme.Set),
				F.Flow2(runInBucket.Get, targetL.Name.Set),
				F.Constant1[runState](targetL.Options.Set(readOptions)),
			}),
		),
		openTarget,
	)
	openOutput = F.Flow2(
		F.Bind2nd(
			build[runState, bucketTarget],
			A.ArrayConcatAll(commonTargets, []func(runState) En.Endomorphism[bucketTarget]{
				F.Flow2(runOutScheme.Get, targetL.Scheme.Set),
				F.Flow2(runOutBucket.Get, targetL.Name.Set),
				F.Constant1[runState](targetL.Options.Set(writeOptions)),
			}),
		),
		openTarget,
	)

	inputRef = F.Bind2nd(build[runState, objectRef], []func(runState) En.Endomorphism[objectRef]{
		F.Flow2(runL.InputBucket.Get, refL.Bucket.Set),
		F.Flow2(runCtx.Get, refL.Ctx.Set),
		F.Flow2(runKey.Get, refL.Key.Set),
	})
	newSource = F.Bind2nd(build[string, source], []func(string) En.Endomorphism[source]{
		sourceL.Key.Set,
		F.Flow2(stemOf, sourceL.Stem.Set),
	})
	memberSource = F.Flow3(
		PR.Of[*zip.File],
		PR.MapHead[*zip.File](F.Flow2(memberName, newSource)),
		PR.Merge(sourceL.Member.Set),
	)
	dirPrefix = F.Flow2(
		F.Bind2nd(strings.TrimSuffix, keySeparator),
		strutils.MapNonEmpty(S.Append(keySeparator)),
	)
	singleSource  = F.Flow4(runKey.Get, newSource, A.Of[source], IOE.Of[error, []source])
	archiveSource = F.Flow4(
		inputRef,
		readRef,
		IOE.ChainEitherK(F.Flow2(openArchive, E.MapLeft[*zip.Reader](ER.OnError("open archive")))),
		IOE.Map[error](F.Flow3(archiveFiles, A.Filter(selectableMember), A.Map(memberSource))),
	)
	dirSource = F.Flow3(
		F.Flow2(inputRef, refL.Key.Modify(dirPrefix)),
		listTwoLevels,
		IOE.Map[error](F.Flow2(A.Filter(selectableObject), A.Map(F.Flow2(objectKey, newSource)))),
	)
	selectSources = MO.Match(
		IOE.Left[[]source](ErrUnknownInputType),
		[]O.Kleisli[runState, IOE.IOEither[error, []source]]{
			MO.Case(F.Flow2(runKind.Get, predord.StrEq(typeDNA)), singleSource),
			MO.Case(F.Flow2(runKind.Get, predord.StrEq(typeZip)), archiveSource),
			MO.Case(F.Flow2(runKind.Get, predord.StrEq(typeDir)), dirSource),
		},
	)
	requireSources = E.FromPredicate(
		F.Flow2(runL.Sources.Get, predarrays.IsNonEmpty[source]()),
		F.Constant1[runState](ErrNoDNAFiles),
	)
	// reserveStems sorts by key and marks every later source of a stem.
	reserveStems = F.Flow4(
		A.SortByKey(S.Ord, sourceL.Key.Get),
		PR.Of[[]source],
		PR.MapTail[[]source](A.MapWithIndex(PR.MakePair[int, source])),
		PR.Paired(markDuplicates),
	)
)

// build applies edits derived from s to the zero value of T.
func build[S, T any](s S, edits []func(S) En.Endomorphism[T]) T {
	return F.Pipe2(edits, A.Flap[En.Endomorphism[T]](s), En.Reduce[T])
}

// update applies edits derived from s to s itself.
func update[S any](s S, edits []func(S) En.Endomorphism[S]) S {
	return F.Pipe3(
		edits,
		A.Flap[En.Endomorphism[S]](s),
		M.ConcatAll(En.Monoid[S]()),
		En.Read(s),
	)
}

func skipSidecars(o fileblob.Options) fileblob.Options {
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
				return fileblob.OpenBucket(dir, opts)
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

func openArchive(data []byte) E.Either[error, *zip.Reader] {
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
			iter := bucket.List(opts)
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

func listOptions(prefix string) *blob.ListOptions {
	return new(F.Pipe2(blob.ListOptions{}, listDelimiter.Set(keySeparator), listPrefix.Set(prefix)))
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

func markDuplicates(sources []source, indexed []PR.Pair[int, source]) []source {
	return F.Pipe1(indexed, A.Map(F.Bind1st(markDuplicate, sources)))
}

func markDuplicate(sources []source, item PR.Pair[int, source]) source {
	return F.Pipe2(
		item,
		PR.MapHead[source](F.Bind1st(earlierSources, sources)),
		PR.Paired(flagDuplicate),
	)
}

func earlierSources(sources []source, index int) []source { return sources[:index] }

func flagDuplicate(earlier []source, src source) source {
	return F.Pipe3(earlier, A.Any(F.Bind1st(sameStem, src)), sourceL.Duplicate.Set, En.Read(src))
}

func sameStem(first, second source) bool {
	return sourceL.Stem.Get(first) == sourceL.Stem.Get(second)
}
