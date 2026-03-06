package logs

import (
	"fmt"
	"io"
	"os"

	fperrors "github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// CreateK8sClient builds a Kubernetes client from the kubeconfig path in Params.
func CreateK8sClient(p Params) IOE.IOEither[error, kubernetes.Interface] {
	return F.Pipe2(
		IOE.TryCatchError(func() (*rest.Config, error) {
			return clientcmd.BuildConfigFromFlags("", p.Kubeconfig)
		}),
		IOE.Chain(func(cfg *rest.Config) IOE.IOEither[error, kubernetes.Interface] {
			return IOE.TryCatchError(func() (kubernetes.Interface, error) {
				return kubernetes.NewForConfig(cfg)
			})
		}),
		IOE.MapLeft[kubernetes.Interface](fperrors.OnError("failed to create k8s client")),
	)
}

// FetchJobPods performs a single K8s API call to list pods matching the job-name label.
// No retry loop — the downstream ExtractLatestPod reports "no pods found" if the list is empty.
func FetchJobPods(ctx LogContext) IOE.IOEither[error, *corev1.PodList] {
	return F.Pipe1(
		IOE.TryCatchError(func() (*corev1.PodList, error) {
			return ctx.Client.CoreV1().Pods(ctx.Namespace).List(
				ctx.Context,
				metav1.ListOptions{LabelSelector: fmt.Sprintf("job-name=%s", ctx.Name)},
			)
		}),
		IOE.MapLeft[*corev1.PodList](fperrors.OnError("failed to list job pods")),
	)
}

// StreamPodLogs opens an HTTP stream to the Kubernetes API to fetch pod logs and writes them to stdout.
func StreamPodLogs(ctx LogContext) IOE.IOEither[error, LogContext] {
	return F.Pipe1(
		IOE.TryCatchError(func() (io.ReadCloser, error) {
			req := ctx.Client.CoreV1().Pods(ctx.Namespace).GetLogs(ctx.PodName, &corev1.PodLogOptions{
				Follow: ctx.Follow,
			})
			return req.Stream(ctx.Context)
		}),
		IOE.Chain(func(stream io.ReadCloser) IOE.IOEither[error, LogContext] {
			return F.Pipe2(
				IOE.TryCatchError(func() (int64, error) {
					defer stream.Close()
					return io.Copy(os.Stdout, stream)
				}),
				IOE.MapLeft[int64](fperrors.OnError("failed to stream logs")),
				IOE.Map[error](func(_ int64) LogContext { return ctx }),
			)
		}),
	)
}
