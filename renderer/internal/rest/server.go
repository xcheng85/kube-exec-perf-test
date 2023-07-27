package rest

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/xcheng85/kube-exec-perf-test/internal/k8s/exec"
	"github.com/xcheng85/kube-exec-perf-test/renderer/internal/service"
	"k8s.io/client-go/kubernetes"
)

// Attach the Renderer service to the server
func RegisterServer(rendererService service.RendererService, ctx context.Context, mux *chi.Mux, exec *exec.K8sExec, clientset kubernetes.Interface) error {
	mux.Post("/renderer", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		e1 := make(chan error)
		go rendererService.Run(exec, clientset, ctx, e1)
		select {
		case <-ctx.Done():
			return
		case <-e1:
			w.Write([]byte(fmt.Sprintf("renderer all done.\n")))
			return
		}
	})
	mux.Post("/long", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		processTime := time.Duration(rand.Intn(4)+1) * time.Second

		select {
		case <-ctx.Done():
			return

		case <-time.After(processTime):
			// The above channel simulates some hard work.
		}

		w.Write([]byte("done"))
	})
	return nil
}
