package gbexport

import (
	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	ER "github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	IO "github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	N "github.com/IBM/fp-go/v2/number"
	L "github.com/IBM/fp-go/v2/optics/lens"
	S "github.com/IBM/fp-go/v2/string"
	predord "github.com/dictyBase/fp-go-loom/predicate/ord"
	"gocloud.dev/blob"
)

var (
	runLogger    = F.Pipe1(runL.State, L.Compose[runState](stateL.Logger))
	runPrimers   = F.Pipe1(runParams, L.Compose[runState](paramsL.ShowPrimers))
	runOutputKey = F.Pipe1(runL.Output, L.Compose[runState](locL.Key))

	// sortByKey orders sources by key so conversion order is stable.
	sortByKey = A.SortByKey(S.Ord, sourceL.Key.Get)
	// requireUniqueStems stops the run when two sources share a stem,
	// before anything is written.
	requireUniqueStems = E.FromPredicate(
		F.Flow2(runL.Sources.Get, hasUniqueStems),
		F.Constant1[runState](ErrDuplicateStem),
	)
)

// Pipeline resolves both locations, opens both buckets with scoped cleanup,
// converts every selected source sequentially (read, parse, encode, write),
// and logs the exported count. The first failure stops the run.
func Pipeline(seed State) IOE.IOEither[error, State] {
	return F.Pipe5(
		runState{},
		runL.State.Set(seed),
		IOE.Of[error, runState],
		IOE.ChainEitherK(resolveLocations),
		IOE.Chain(withInputBucket),
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

// withInputBucket opens the input bucket and closes it on success or failure.
func withInputBucket(r runState) IOE.IOEither[error, State] {
	return F.Pipe2(
		r.withOutputBucket,
		IOE.WithResource[State](openInputBucket(r), closeResource[*blob.Bucket]),
		IOE.MapLeft[State](ER.OnError("export")),
	)
}

// withOutputBucket opens the output bucket with its own resource scope,
// converts all sources, and closes the output bucket on success or failure.
func (r runState) withOutputBucket(in *blob.Bucket) IOE.IOEither[error, State] {
	return F.Pipe1(
		F.Bind1st(r.convertAll, in),
		IOE.WithResource[State](openOutputBucket(r), closeResource[*blob.Bucket]),
	)
}

// convertAll selects sources from the input bucket, validates them, and
// converts them sequentially into the output bucket.
func (r runState) convertAll(in, out *blob.Bucket) IOE.IOEither[error, State] {
	return F.Pipe8(
		r,
		runL.InputBucket.Set(in),
		runL.OutputBucket.Set(out),
		IOE.Of[error, runState],
		IOE.Bind(runL.Sources.Set, selectSources),
		IOE.ChainEitherK(requireSources),
		IOE.ChainEitherK(requireUniqueStems),
		IOE.Map[error](runL.Sources.Modify(sortByKey)),
		IOE.Chain(convertSources),
	)
}

// convertSources converts sources sequentially and counts successes.
// A source starts only after the previous one succeeded; IOE.TraverseArraySeq
// is avoided because it runs every effect before combining results.
func convertSources(r runState) IOE.IOEither[error, State] {
	return F.Pipe4(
		r,
		runL.Sources.Get,
		A.Map(F.Bind1st(jobFor, r)),
		A.Reduce(convertNext, IOE.Of[error](0)),
		IOE.Map[error](r.exported),
	)
}

// convertNext converts job only if all earlier jobs succeeded, then increments count.
func convertNext(done IOE.IOEither[error, int], job fileJob) IOE.IOEither[error, int] {
	effect := convertSource(job)
	return F.Pipe2(
		done,
		IOE.ChainFirst(F.Constant1[int](effect)),
		IOE.Map[error](N.Add(1)),
	)
}

// exported returns the state of r with count recorded as Exported.
func (r runState) exported(count int) State {
	return F.Pipe2(r, runL.State.Get, stateL.Exported.Set(count))
}

// jobFor constructs the conversion job for src in run r.
func jobFor(r runState, src source) fileJob {
	return F.Pipe2(
		fileJob{},
		jobL.Run.Set(r),
		jobL.Source.Set(src),
	)
}

// hasUniqueStems reports whether every source has a distinct stem.
func hasUniqueStems(sources []source) bool {
	total := len(sources)
	return F.Pipe3(
		sources,
		A.Uniq(sourceL.Stem.Get),
		A.Size[source],
		F.Bind1st(predord.IntEq.Equals, total),
	)
}
