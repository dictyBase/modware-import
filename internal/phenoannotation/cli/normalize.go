// Package cli holds the actions of the phenoannotation commands.
package cli

import (
	"context"
	"errors"

	E "github.com/IBM/fp-go/v2/either"
	En "github.com/IBM/fp-go/v2/endomorphism"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/predicate"
	endoutils "github.com/dictyBase/fp-go-loom/endomorphism"
	ioeutils "github.com/dictyBase/fp-go-loom/ioeitherutils"
	MO "github.com/dictyBase/fp-go-loom/matchopt"
	cli "github.com/urfave/cli/v3"

	norm "github.com/dictyBase/modware-import/internal/phenoannotation/normalize"
)

// Exit codes of the process.
const (
	exitOK     = 0
	exitFailed = 1
	exitReview = 2
)

var (
	// ErrReview marks a complete batch that contains rejected rows.
	ErrReview = errors.New("batch contains rejected rows: review rejects.csv")

	stateL   = norm.MakeStateLenses()
	optionsL = norm.MakeOptionsLenses()
	batchL   = norm.MakeBatchLenses()

	stateRejected = F.Pipe1(stateL.Batch, L.Compose[norm.State](batchL.Rejected))

	isReview = F.Bind2nd(errors.Is, ErrReview)

	// ExitCode maps a command error to the process exit code: 0 for no
	// error, 2 for a published batch with rejected rows, and 1 otherwise.
	ExitCode = MO.Match(exitFailed, []O.Kleisli[error, int]{
		MO.Const(P.IsZero[error](), exitOK),
		MO.Const(isReview, exitReview),
	})

	// requireNoRejects turns a batch with rejected rows into ErrReview.
	// The batch is already published when this step runs.
	requireNoRejects = E.FromPredicate(
		F.Flow2(stateRejected.Get, P.IsZero[int]()),
		F.Constant1[norm.State](ErrReview),
	)

	flagText = (*cli.Command).String

	// toOptions copies the flag values into the options.
	toOptions = endoutils.BuildFrom([]func(*cli.Command) En.Endomorphism[norm.Options]{
		F.Flow2(F.Bind2nd(flagText, "folder"), optionsL.Folder.Set),
		F.Flow2(F.Bind2nd(flagText, "sheet"), optionsL.Sheet.Set),
		F.Flow2(F.Bind2nd(flagText, "out"), optionsL.Out.Set),
	})
)

// RunNormalize is the action of the normalize subcommand. It validates the
// flags, selects files, reads and classifies every row, publishes the
// batch, prints the summary, and returns the review status. It forces the
// effect once.
func RunNormalize(ctx context.Context, c *cli.Command) error {
	return F.Pipe11(
		norm.State{},
		IOE.Of[error, norm.State],
		IOE.LetTo[error](stateL.Context.Set, ctx),
		IOE.ApS(stateL.Options.Set, F.Pipe2(c, toOptions, IOE.Of[error, norm.Options])),
		IOE.ChainEitherK(norm.Valid),
		IOE.Bind(stateL.Files.Set, norm.SelectFiles),
		IOE.Bind(stateL.Batch.Set, norm.ReadFiles),
		IOE.Chain(norm.Publish),
		IOE.ChainFirstIOK[error](norm.PrintSummary),
		IOE.ChainEitherK(requireNoRejects),
		ioeutils.ToEither[error, norm.State],
		E.ToError[norm.State],
	)
}
