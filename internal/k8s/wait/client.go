package wait

import (
	"context"
	"fmt"

	fperrors "github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// FetchJob retrieves the named Job from Kubernetes and stores it in ctx.Job.
func FetchJob(ctx PollContext) IOE.IOEither[error, PollContext] {
	return F.Pipe2(
		IOE.TryCatchError(func() (*batchv1.Job, error) {
			return ctx.Client.BatchV1().Jobs(ctx.Namespace).Get(
				context.Background(), ctx.Name, metav1.GetOptions{},
			)
		}),
		IOE.MapLeft[*batchv1.Job](fperrors.OnError("failed to get job")),
		IOE.Map[error](func(job *batchv1.Job) PollContext {
			ctx.Job = job
			return ctx
		}),
	)
}

// FetchPods lists pods belonging to the named Job.
func FetchPods(ctx PollContext) IOE.IOEither[error, *corev1.PodList] {
	return F.Pipe1(
		IOE.TryCatchError(func() (*corev1.PodList, error) {
			return ctx.Client.CoreV1().Pods(ctx.Namespace).List(
				context.Background(),
				metav1.ListOptions{
					LabelSelector: fmt.Sprintf("job-name=%s", ctx.Name),
				},
			)
		}),
		IOE.MapLeft[*corev1.PodList](
			fperrors.OnError("failed to list job pods"),
		),
	)
}

// CreateK8sClient builds a Kubernetes client from the kubeconfig path in Params.
// Falls back to in-cluster config when Kubeconfig is empty.
func CreateK8sClient(p Params) IOE.IOEither[error, kubernetes.Interface] {
	return F.Pipe2(
		IOE.TryCatchError(func() (*rest.Config, error) {
			return clientcmd.BuildConfigFromFlags("", p.Kubeconfig)
		}),
		IOE.Chain(
			func(cfg *rest.Config) IOE.IOEither[error, kubernetes.Interface] {
				return IOE.TryCatchError(func() (kubernetes.Interface, error) {
					return kubernetes.NewForConfig(cfg)
				})
			},
		),
		IOE.MapLeft[kubernetes.Interface](
			fperrors.OnError("failed to create k8s client"),
		),
	)
}
