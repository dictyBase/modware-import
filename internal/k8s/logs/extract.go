package logs

import (
        "fmt"
        "slices"

        A "github.com/IBM/fp-go/v2/array"
        F "github.com/IBM/fp-go/v2/function"
        IOE "github.com/IBM/fp-go/v2/ioeither"
        corev1 "k8s.io/api/core/v1"
)

// sortPodsDesc clones and sorts a slice of pods by their CreationTimestamp in descending order (newest first).
func sortPodsDesc(pods []corev1.Pod) []corev1.Pod {
        sorted := slices.Clone(pods)
        slices.SortFunc(sorted, func(a, b corev1.Pod) int {
                // Return -1 if b is before a, meaning newer items bubble to the top
                return b.CreationTimestamp.Time.Compare(a.CreationTimestamp.Time)
        })
        return sorted
}

// ExtractLatestPod uses functional pipelines to sort and safely retrieve the newest pod from the LogContext.
func ExtractLatestPod(ctx LogContext) IOE.IOEither[error, LogContext] {
        return F.Pipe4(
                ctx.Pods.Items,
                sortPodsDesc,
                A.Head[corev1.Pod],
                IOE.FromOption[corev1.Pod](func() error {
                        return fmt.Errorf("no pods found for job")
                }),
                IOE.Map[error](func(pod corev1.Pod) LogContext {
                        ctx.PodName = pod.Name
                        return ctx
                }),
        )
}
