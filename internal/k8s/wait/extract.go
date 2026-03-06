package wait

import (
	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

// stuckReasons is the set of container waiting reasons that indicate a stuck pod.
var stuckReasons = map[string]bool{
	"ImagePullBackOff":           true,
	"ErrImagePull":               true,
	"InvalidImageName":           true,
	"CrashLoopBackOff":           true,
	"CreateContainerConfigError": true,
	"CreateContainerError":       true,
	"ContainerCannotRun":         true,
}

// isStuckReason returns true if the waiting reason signals a stuck container.
func isStuckReason(reason string) bool {
	return stuckReasons[reason]
}

// allContainerStatuses concatenates init and regular container statuses for a pod.
func allContainerStatuses(pod corev1.Pod) []corev1.ContainerStatus {
	return append(
		pod.Status.InitContainerStatuses,
		pod.Status.ContainerStatuses...)
}

// toStuckReason extracts a stuck reason from a ContainerStatus if one exists.
// Uses O.FromNillable to safely handle the nullable Waiting pointer.
func toStuckReason(cs corev1.ContainerStatus) O.Option[string] {
	return F.Pipe3(
		cs.State.Waiting,
		O.FromNillable,
		O.Map(func(w *corev1.ContainerStateWaiting) string { return w.Reason }),
		O.Chain(O.FromPredicate(isStuckReason)),
	)
}

// FindStuckReason scans all pods and all containers for a known stuck reason.
// Returns Some(reason) on first match, None if all containers are healthy.
func FindStuckReason(pods *corev1.PodList) O.Option[string] {
	return F.Pipe3(
		pods.Items,
		A.Chain(allContainerStatuses),
		A.FilterMap(toStuckReason),
		A.Head,
	)
}

// ClassifyPodState returns JobStuck if any container is stuck, otherwise JobPending.
func ClassifyPodState(pods *corev1.PodList) JobState {
	return F.Pipe1(
		FindStuckReason(pods),
		O.Fold(
			F.Constant(JobPending),
			func(_ string) JobState { return JobStuck },
		),
	)
}

// conditionTypeMap maps Kubernetes JobConditionTypes to domain JobStates.
var conditionTypeMap = map[batchv1.JobConditionType]JobState{
	batchv1.JobComplete: JobComplete,
	batchv1.JobFailed:   JobFailed,
}

// isTerminalCondition returns true for a Complete or Failed condition with Status=True.
func isTerminalCondition(c batchv1.JobCondition) bool {
	return c.Status == corev1.ConditionTrue &&
		(c.Type == batchv1.JobComplete || c.Type == batchv1.JobFailed)
}

// conditionToJobState maps a terminal JobCondition to the corresponding JobState.
func conditionToJobState(c batchv1.JobCondition) JobState {
	return conditionTypeMap[c.Type]
}

// toTerminalState converts a single JobCondition into an Option[JobState].
// Returns None for non-terminal or Status!=True conditions.
func toTerminalState(c batchv1.JobCondition) O.Option[JobState] {
	return F.Pipe2(
		c,
		O.FromPredicate(isTerminalCondition),
		O.Map(conditionToJobState),
	)
}

// ExtractJobCondition reads ctx.Job, extracts the first terminal condition,
// and stores it in ctx.Condition.
func ExtractJobCondition(ctx PollContext) PollContext {
	return F.Pipe3(
		ctx.Job.Status.Conditions,
		A.FilterMap(toTerminalState),
		A.Head,
		func(cond O.Option[JobState]) PollContext {
			ctx.Condition = cond
			return ctx
		},
	)
}
