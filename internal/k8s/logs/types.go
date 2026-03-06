package logs

import (
        "context"

        F "github.com/IBM/fp-go/v2/function"
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
}

var SetClient = F.Curry2(func(c kubernetes.Interface, p Params) LogContext {
        return LogContext{Params: p, Client: c}
})
