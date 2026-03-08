package wait

import (
	"log/slog"
	"time"

	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/client-go/kubernetes"
)

// JobState represents the observed state of a Kubernetes Job.
type JobState string

const (
	JobPending  JobState = "Pending"
	JobComplete JobState = "Complete"
	JobFailed   JobState = "Failed"
	JobStuck    JobState = "Stuck"
)

// PollInterval is the delay between polling cycles.
const PollInterval = 5 * time.Second

// Params holds CLI-supplied parameters for the wait-job command.
type Params struct {
	Name       string
	Namespace  string
	Timeout    time.Duration
	Kubeconfig string
	Logger     *slog.Logger
}

// WithClient enriches Params with an injected Kubernetes client.
type WithClient struct {
	Params
	Client kubernetes.Interface
}

// PollContext is the fully-enriched context used throughout the polling loop.
type PollContext struct {
	WithClient
	Deadline  time.Time
	Job       *batchv1.Job
	Condition O.Option[JobState]
	State     JobState
}

// SetClient is a curried setter used with IOE.Bind to inject the K8s client.
var SetClient = F.Curry2(func(c kubernetes.Interface, p Params) WithClient {
	return WithClient{Params: p, Client: c}
})

// SetPollReady is a curried setter used with IOE.Let to attach the deadline.
// The Logger is preserved from Params.
var SetPollReady = F.Curry2(func(deadline time.Time, c WithClient) PollContext {
	return PollContext{
		WithClient: c,
		Deadline:   deadline,
		State:      JobPending,
		Condition:  O.None[JobState](),
	}
})

// computeDeadline derives the polling deadline from Params.Timeout.
// Pure function — safe for use with IOE.Let.
func computeDeadline(c WithClient) time.Time {
	return time.Now().Add(c.Timeout)
}

// isTerminal returns true for any JobState that ends the polling loop.
func isTerminal(s JobState) bool {
	return s == JobComplete || s == JobFailed || s == JobStuck
}

// isComplete returns true only for a successfully completed job.
func isComplete(s JobState) bool {
	return s == JobComplete
}
