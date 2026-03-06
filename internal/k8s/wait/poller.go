package wait

import (
	"fmt"
	"time"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	O "github.com/IBM/fp-go/v2/option"
)

// checkTimeout fails with a timeout error if the deadline has passed.
// Passes PollContext forward on success so FetchJob can be chained point-free.
func checkTimeout(ctx PollContext) E.Either[error, PollContext] {
	return F.Pipe1(
		ctx,
		E.FromPredicate(
			func(c PollContext) bool {
				return time.Now().Before(c.Deadline)
			},
			func(c PollContext) error {
				return fmt.Errorf("timeout waiting for job %s", c.Name)
			},
		),
	)
}

// resolveState reads ctx.Condition and either uses it directly (terminal) or
// fetches pods to classify (Pending vs Stuck), updating ctx.State.
func resolveState(ctx PollContext) IOE.IOEither[error, PollContext] {
	withState := func(state JobState) PollContext {
		ctx.State = state
		return ctx
	}
	return F.Pipe1(
		ctx.Condition,
		O.Fold(
			func() IOE.IOEither[error, PollContext] {
				return F.Pipe3(
					ctx,
					FetchPods,
					IOE.Map[error](ClassifyPodState),
					IOE.Map[error](withState),
				)
			},
			func(state JobState) IOE.IOEither[error, PollContext] {
				return IOE.Of[error](withState(state))
			},
		),
	)
}

// continueOrReturn decides whether to stop or keep polling based on ctx.State.
// Terminal state (Complete/Failed/Stuck) returns ctx immediately.
// Pending sleeps for PollInterval and recurses directly into pollUntilDone.
func continueOrReturn(ctx PollContext) IOE.IOEither[error, PollContext] {
	return F.Pipe2(
		ctx.State,
		O.FromPredicate(isTerminal),
		O.Fold(
			func() IOE.IOEither[error, PollContext] {
				time.Sleep(PollInterval)
				return pollUntilDone(ctx)
			},
			func(_ JobState) IOE.IOEither[error, PollContext] {
				return IOE.Of[error](ctx)
			},
		),
	)
}

// pollUntilDone is the recursive polling loop. Each iteration: check timeout →
// fetch condition → resolve state → log → continue or return. All pipe steps
// are point-free — PollContext threads through without outer closures.
func pollUntilDone(ctx PollContext) IOE.IOEither[error, PollContext] {
	return F.Pipe7(
		ctx,
		checkTimeout,
		IOE.FromEither[error, PollContext],
		IOE.Chain(FetchJob),
		IOE.Map[error](ExtractJobCondition),
		IOE.Chain(resolveState),
		IOE.ChainFirstIOK[error](logState),
		IOE.Chain(continueOrReturn),
	)
}
