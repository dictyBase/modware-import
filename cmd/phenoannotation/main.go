// Command phenoannotation prepares phenotype annotation batches.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	F "github.com/IBM/fp-go/v2/function"
	IO "github.com/IBM/fp-go/v2/io"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/predicate"
	MO "github.com/dictyBase/fp-go-loom/matchopt"

	app "github.com/dictyBase/modware-import/internal/phenoannotation/cli"
)

// Exit codes of the process.
const (
	exitOK     = 0
	exitFailed = 1
	exitReview = 2
)

var (
	isReview = F.Bind2nd(errors.Is, app.ErrReview)

	// exitCode maps the command error to the process exit code.
	exitCode = MO.Match(exitFailed, []O.Kleisli[error, int]{
		MO.Const(P.IsZero[error](), exitOK),
		MO.Const(isReview, exitReview),
	})

	// reportFailure describes the stderr report of a command error. It keeps
	// the error and does nothing for a nil error.
	reportFailure = F.Pipe1(P.IsNonZero[error](), P.Fold(IO.Of[error], printFailure))
)

// printFailure returns the effect that writes err to stderr and keeps err.
func printFailure(err error) IO.IO[error] {
	return F.Pipe1(
		IO.FromImpure(func() { _, _ = fmt.Fprintln(os.Stderr, err) }),
		IO.MapTo[IO.Void](err),
	)
}

// main runs the command. A command action never calls os.Exit, so the
// cleanup of the action runs before main stops the signal handler.
func main() {
	background := context.Background()
	ctx, stop := signal.NotifyContext(background, os.Interrupt, syscall.SIGTERM)
	err := app.Command(os.Stdout).Run(ctx, os.Args)
	stop()
	os.Exit(F.Pipe3(err, reportFailure, IO.Map(exitCode), IO.Run))
}
