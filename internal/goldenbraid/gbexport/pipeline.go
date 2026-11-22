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
	endoutils "github.com/dictyBase/fp-go-loom/endomorphism"
	"gocloud.dev/blob"

	"github.com/dictyBase/modware-import/internal/goldenbraid/seqmodel"
)

var (
	summaryL     = seqmodel.MakeExportSummaryLenses()
	stateTotal   = F.Pipe1(stateL.Summary, L.Compose[State](summaryL.Total))
	stateDone    = F.Pipe1(stateL.Summary, L.Compose[State](summaryL.Succeeded))
	stateFailed  = F.Pipe1(stateL.Summary, L.Compose[State](summaryL.Failed))
	stateErrors  = F.Pipe1(stateL.Summary, L.Compose[State](summaryL.Errors))
	runLogger    = F.Pipe1(runL.State, L.Compose[runState](stateL.Logger))
	runPrimers   = F.Pipe1(runParams, L.Compose[runState](paramsL.ShowPrimers))
	runOutputKey = F.Pipe1(runL.Output, L.Compose[runState](locL.Key))

	failOnPartial = E.FromPredicate(F.Flow2(stateFailed.Get, isZero), partialFailure)
	jobEdits      = []func(PR.Pair[runState, source]) En.Endomorphism[fileJob]{
		F.Flow3(PR.Head[runState, source], runCtx.Get, jobL.Ctx.Set),
		F.Flow3(PR.Head[runState, source], runLogger.Get, jobL.Logger.Set),
		F.Flow3(PR.Head[runState, source], runPrimers.Get, jobL.ShowPrimers.Set),
		F.Flow3(PR.Head[runState, source], runL.InputBucket.Get, jobL.Input.Set),
		F.Flow3(PR.Head[runState, source], runL.OutputBucket.Get, jobL.Output.Set),
		F.Flow3(PR.Head[runState, source], runOutputKey.Get, jobL.Prefix.Set),
		F.Flow2(PR.Tail[runState, source], jobL.Source.Set),
	}
	// jobOf makes the export job of one source of a run.
	jobOf = endoutils.BuildFrom(jobEdits)
)

// Pipeline resolves both locations, exports every selected construct one at
// a time inside the bucket lifetimes, logs the counts once, and returns
// ErrPartialExportFailure when any file failed.
func Pipeline(seed State) IOE.IOEither[error, State] {
	return F.Pipe6(
		runState{},
		runL.State.Set(seed),
		IOE.Of[error, runState],
		IOE.ChainEitherK(resolveLocations),
		IOE.Chain(exportInput),
		IOE.ChainFirstIOK[error](State.logSummary),
		IOE.ChainEitherK(failOnPartial),
	)
}

func (s State) logSummary() IO.IO[IO.Void] {
	return IO.FromImpure(
		func() {
			stateL.Logger.Get(s).Info(
				"export summary",
				"total", stateTotal.Get(s),
				"succeeded", stateDone.Get(s),
				"failed", stateFailed.Get(s),
			)
		},
	)
}

// exportInput keeps the input bucket open while it selects the sources and
// runs the output scope.
func exportInput(r runState) IOE.IOEither[error, State] {
	return F.Pipe3(
		r,
		openInput,
		F.Bind2nd(withResource[*blob.Bucket, State], F.Bind1st(withInputBucket, r)),
		IOE.MapLeft[State](ER.OnError("export input")),
	)
}

func withInputBucket(r runState, bucket *blob.Bucket) IOE.IOEither[error, State] {
	return F.Pipe6(
		r,
		runL.InputBucket.Set(bucket),
		IOE.Of[error, runState],
		IOE.Bind(runL.Sources.Set, selectSources),
		IOE.ChainEitherK(requireSources),
		IOE.Map[error](runL.Sources.Modify(reserveStems)),
		IOE.Chain(exportOutput),
	)
}

// exportOutput opens the output only after the input has sources.
func exportOutput(r runState) IOE.IOEither[error, State] {
	return F.Pipe3(
		r,
		openOutput,
		F.Bind2nd(withResource[*blob.Bucket, State], F.Bind1st(withOutputBucket, r)),
		IOE.MapLeft[State](ER.OnError("export output")),
	)
}

// withOutputBucket exports the sources sequentially and keeps each failure
// in the summary rather than stopping the batch.
func withOutputBucket(r runState, bucket *blob.Bucket) IOE.IOEither[error, State] {
	return F.Pipe6(
		r,
		runL.OutputBucket.Set(bucket),
		jobsOf,
		IO.TraverseArraySeq(exportFile),
		IO.Map(seqmodel.ReduceExportResults),
		IOE.FromIO[error, seqmodel.ExportSummary],
		IOE.Map[error](F.Pipe2(r, runL.State.Get, F.Flip(stateL.Summary.Set))),
	)
}

// jobsOf makes one export job for every source of r.
func jobsOf(r runState) []fileJob {
	return F.Pipe2(r, runL.Sources.Get, A.Map(F.Bind1st(jobFor, r)))
}

// jobFor makes the export job of src in run r.
func jobFor(r runState, src source) fileJob {
	return F.Pipe2(src, F.Bind1st(PR.MakePair[runState, source], r), jobOf)
}

func isZero(n int) bool { return n == 0 }

// partialFailure keeps every per-file cause available to errors.Is.
func partialFailure(s State) error {
	causes := stateErrors.Get(s)
	return errors.Join(ErrPartialExportFailure, errors.Join(causes...))
}
