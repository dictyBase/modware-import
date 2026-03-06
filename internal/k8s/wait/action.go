package wait

import (
        "fmt"
        "log/slog"
        "time"

        E "github.com/IBM/fp-go/v2/either"
        F "github.com/IBM/fp-go/v2/function"
        IOE "github.com/IBM/fp-go/v2/ioeither"
        "github.com/dictyBase/modware-import/internal/logger"
        "github.com/urfave/cli/v2"
)

const defaultTimeout = 60 * time.Second

// validateTerminalState succeeds for Complete; fails with an error for any other state.
var validateTerminalState = E.FromPredicate(
        isComplete,
        func(s JobState) error {
                return fmt.Errorf("job terminated in state: %s", s)
        },
)

// parseDuration converts a CLI timeout string to time.Duration.
// Falls back to defaultTimeout (60s) on parse failure.
func parseDuration(s string) time.Duration {
        return F.Pipe1(
                E.TryCatchError(time.ParseDuration(s)),
                E.GetOrElse(func(_ error) time.Duration { return defaultTimeout }),
        )
}

// extractState pulls the final JobState out of a completed PollContext.
func extractState(ctx PollContext) JobState { return ctx.State }

// toEither executes an IOEither to get an Either result.
func toEither[ER, A any](ioe IOE.IOEither[ER, A]) E.Either[ER, A] {
        return ioe()
}

// JobAction is the urfave/cli v2 action for the wait-job subcommand.
//
// Pipeline (8 steps):
//  1. Build Params from CLI flags
//  2. Inject Kubernetes client (Bind)
//  3. Compute deadline and attach logger (Let)
//  4. Run polling loop (returns PollContext with final State)
//  5. Extract JobState from PollContext
//  6. Execute IOEither effect → Either
//  7. Validate terminal state (Complete → ok, else → error)
//  8. Fold Either → error
func JobAction(c *cli.Context) error {
        logHandler := logger.GetCliSlogHandler(c)
        slogger := slog.New(logHandler)

        return F.Pipe7(
                IOE.Of[error](Params{
                        Name:       c.String("name"),
                        Namespace:  c.String("namespace"),
                        Timeout:    parseDuration(c.String("timeout")),
                        Kubeconfig: c.String("kubeconfig"),
                        Logger:     slogger,
                }),
                IOE.Bind(SetClient, CreateK8sClient),
                IOE.Let[error](SetPollReady, computeDeadline),
                IOE.Chain(pollUntilDone),
                IOE.Map[error](extractState),
                toEither[error, JobState],
                E.Chain(validateTerminalState),
                E.Fold(
                        F.Identity[error],
                        func(_ JobState) error { return nil },
                ),
        )
}
