package gbexport

import (
	"errors"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	En "github.com/IBM/fp-go/v2/endomorphism"
	ER "github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	IO "github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	L "github.com/IBM/fp-go/v2/optics/lens"
	PR "github.com/IBM/fp-go/v2/pair"
	P "github.com/IBM/fp-go/v2/predicate"
	S "github.com/IBM/fp-go/v2/string"
	endoutils "github.com/dictyBase/fp-go-loom/endomorphism"
	"gocloud.dev/blob"
)

var (
	bucketsL     = MakebucketsLenses()
	runInBlob    = F.Pipe1(runL.Buckets, L.Compose[runState](bucketsL.In))
	runOutBlob   = F.Pipe1(runL.Buckets, L.Compose[runState](bucketsL.Out))
	runLogger    = F.Pipe1(runL.State, L.Compose[runState](stateL.Logger))
	runPrimers   = F.Pipe1(runParams, L.Compose[runState](paramsL.ShowPrimers))
	runOutputKey = F.Pipe1(runL.Output, L.Compose[runState](locL.Key))

	// sortByKey orders the sources by key so the conversion order is stable.
	sortByKey = A.SortByKey(S.Ord, sourceL.Key.Get)
	// requireUniqueStems stops the run when two sources share a stem,
	// before anything is written.
	requireUniqueStems = E.FromPredicate(
		F.Flow3(runL.Sources.Get, duplicateStems, P.IsZero[int]()),
		F.Constant1[runState](ErrDuplicateStem),
	)
	jobEdits = []func(PR.Pair[runState, source]) En.Endomorphism[fileJob]{
		F.Flow3(PR.Head[runState, source], runCtx.Get, jobL.Ctx.Set),
		F.Flow3(PR.Head[runState, source], runLogger.Get, jobL.Logger.Set),
		F.Flow3(PR.Head[runState, source], runPrimers.Get, jobL.ShowPrimers.Set),
		F.Flow3(PR.Head[runState, source], runInBlob.Get, jobL.Input.Set),
		F.Flow3(PR.Head[runState, source], runOutBlob.Get, jobL.Output.Set),
		F.Flow3(PR.Head[runState, source], runOutputKey.Get, jobL.Prefix.Set),
		F.Flow2(PR.Tail[runState, source], jobL.Source.Set),
	}
	// jobOf makes the conversion job of one source of a run.
	jobOf = endoutils.BuildFrom(jobEdits)
)

// Pipeline resolves both locations, opens both buckets, converts every
// selected source one at a time (read, parse, encode, write), closes
// both buckets, and logs the count. The first failure stops the run.
func Pipeline(seed State) IOE.IOEither[error, State] {
	return F.Pipe5(
		runState{},
		runL.State.Set(seed),
		IOE.Of[error, runState],
		IOE.ChainEitherK(resolveLocations),
		IOE.Chain(copySources),
		IOE.ChainFirstIOK[error](State.logSummary),
	)
}

func (s State) logSummary() IO.IO[IO.Void] {
	return IO.FromImpure(
		func() {
			stateL.Logger.Get(s).Info("export summary", "exported", stateL.Exported.Get(s))
		},
	)
}

// copySources opens both buckets, converts every source of the input
// into the output, and closes both buckets on success and on failure.
func copySources(r runState) IOE.IOEither[error, State] {
	return F.Pipe2(
		r.convertAll,
		IOE.WithResource[State](openBuckets(r), closeBuckets),
		IOE.MapLeft[State](ER.OnError("export")),
	)
}

// convertAll selects the sources in the input bucket and converts them
// one at a time into the output bucket. It stops at the first failure.
func (r runState) convertAll(open buckets) IOE.IOEither[error, State] {
	return F.Pipe7(
		r,
		runL.Buckets.Set(open),
		IOE.Of[error, runState],
		IOE.Bind(runL.Sources.Set, selectSources),
		IOE.ChainEitherK(requireSources),
		IOE.ChainEitherK(requireUniqueStems),
		IOE.Map[error](runL.Sources.Modify(sortByKey)),
		IOE.Chain(convertSources),
	)
}

// convertSources converts the sources of r in order and returns the
// state with the number of exported sources. A source starts only
// after the previous one succeeded. IOE.TraverseArraySeq is not used
// because it runs every effect before it combines the results.
func convertSources(r runState) IOE.IOEither[error, State] {
	return F.Pipe3(
		r,
		jobsOf,
		A.Reduce(convertNext, IOE.Of[error](0)),
		IOE.Map[error](r.exported),
	)
}

// convertNext converts job when every earlier job succeeded.
func convertNext(done IOE.IOEither[error, int], job fileJob) IOE.IOEither[error, int] {
	return F.Pipe1(done, IOE.Chain(F.Bind2nd(convertAndCount, job)))
}

// convertAndCount converts job and adds one to count.
func convertAndCount(count int, job fileJob) IOE.IOEither[error, int] {
	return F.Pipe2(job, convertSource, IOE.Map[error](F.Constant1[string](count+1)))
}

// exported returns the state of r with count as the exported sources.
func (r runState) exported(count int) State {
	return F.Pipe2(r, runL.State.Get, stateL.Exported.Set(count))
}

// jobsOf makes one conversion job for every source of r.
func jobsOf(r runState) []fileJob {
	return F.Pipe2(r, runL.Sources.Get, A.Map(F.Bind1st(jobFor, r)))
}

// jobFor makes the conversion job of src in run r.
func jobFor(r runState, src source) fileJob {
	return F.Pipe2(src, F.Bind1st(PR.MakePair[runState, source], r), jobOf)
}

// duplicateStems counts the sources whose stem an earlier source
// already has.
func duplicateStems(sources []source) int {
	return len(sources) - F.Pipe2(sources, A.Uniq(sourceL.Stem.Get), A.Size[source])
}

// openBuckets opens the input bucket, then the output bucket. It closes
// the input bucket when the output bucket cannot open.
func openBuckets(r runState) IOE.IOEither[error, buckets] {
	return F.Pipe1(openInput(r), IOE.Chain(r.withOutput))
}

// withOutput opens the output bucket next to the open input bucket in.
func (r runState) withOutput(in *blob.Bucket) IOE.IOEither[error, buckets] {
	return F.Pipe2(
		openOutput(r),
		IOE.Map[error](F.Bind1st(newBuckets, in)),
		IOE.OrElse(F.Bind1st(closeAfterFailure, in)),
	)
}

func newBuckets(in, out *blob.Bucket) buckets {
	return F.Pipe2(buckets{}, bucketsL.In.Set(in), bucketsL.Out.Set(out))
}

// closeAfterFailure closes in and returns err. A close error joins err.
func closeAfterFailure(in *blob.Bucket, err error) IOE.IOEither[error, buckets] {
	return F.Pipe2(
		closeResource(in),
		IOE.MapLeft[struct{}](F.Bind1st(joinErrors, err)),
		IOE.Chain(F.Constant1[struct{}](IOE.Left[buckets](err))),
	)
}

func joinErrors(first, second error) error { return errors.Join(first, second) }

// closeBuckets closes both buckets and joins their close errors.
func closeBuckets(open buckets) IOE.IOEither[error, struct{}] {
	return IOE.TryCatchError(
		func() (struct{}, error) {
			in, out := bucketsL.In.Get(open), bucketsL.Out.Get(open)
			return struct{}{}, errors.Join(in.Close(), out.Close())
		},
	)
}
