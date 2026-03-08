package logs

import (
	"context"

	F "github.com/IBM/fp-go/v2/function"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type Params struct {
	Name       string
	Namespace  string
	Kubeconfig string
	Follow     bool
	Context    context.Context
}

type LogContext struct {
	Params
	Client  kubernetes.Interface
	PodName string
	Pods    *corev1.PodList
}

var SetClient = F.Curry2(func(c kubernetes.Interface, p Params) LogContext {
	return LogContext{Params: p, Client: c}
})
