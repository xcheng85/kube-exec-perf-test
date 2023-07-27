package monolith

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/xcheng85/kube-exec-perf-test/internal/k8s/exec"
	"k8s.io/client-go/kubernetes"
)

// define the interface of di
// Module requires Monolith interface

// chi.Mux is the implementation of chi.Router interface
type Monolith interface {
	Mux() *chi.Mux
	Exec() *exec.K8sExec
	KubernetesClientSet() kubernetes.Interface
}

type Module interface {
	Startup(context.Context, Monolith) error
}
