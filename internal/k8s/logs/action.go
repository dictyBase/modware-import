package logs

import (
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	"github.com/urfave/cli/v2"
	corev1 "k8s.io/api/core/v1"
)

func toEither[ER, A any](
	ioe IOE.IOEither[ER, A],
) E.Either[ER, A] {
	return ioe()
}

// JobLogsAction is the urfave/cli v2 action for the job-logs subcommand.
//
// Pipeline (6 steps):
//  1. Build Params from CLI flags (including context for cancellation)
//  2. Inject Kubernetes client (Bind)
//  3. Fetch pods for the job via single API call
//  4. Extract the most recently created pod from the list
//  5. Stream the logs from that pod directly to stdout
//  6. Fold the Either into a standard Go error
func JobLogsAction(c *cli.Context) error {
	return F.Pipe4(
		IOE.Of[error](Params{
			Name:       c.String("name"),
			Namespace:  c.String("namespace"),
			Kubeconfig: c.String("kubeconfig"),
			Follow:     c.Bool("follow"),
			Context:    c.Context,
		}),
		IOE.Bind(SetClient, CreateK8sClient),
		IOE.Chain(func(ctx LogContext) IOE.IOEither[error, LogContext] {
			return F.Pipe4(
				ctx,
				FetchJobPods,
				IOE.Chain(ExtractLatestPod),
				IOE.Map[error](func(pod corev1.Pod) LogContext {
					ctx.PodName = pod.Name
					return ctx
				}),
				IOE.Chain(StreamPodLogs),
			)
		}),
		toEither[error, LogContext],
		E.Fold(
			F.Identity[error],
			func(_ LogContext) error { return nil },
		),
	)
}
