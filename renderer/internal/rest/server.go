package rest

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/xcheng85/kube-exec-perf-test/internal/k8s/exec"
	"github.com/xcheng85/kube-exec-perf-test/renderer/internal/service"
	"net/http"
)

// Attach the Renderer service to the server
func RegisterServer(rendererService service.RendererService, ctx context.Context, mux *chi.Mux, exec *exec.K8sExec) error {
	mux.Post("/renderer", func(w http.ResponseWriter, r *http.Request) {
		rendererService.Run(exec, ctx)
		w.Write([]byte(fmt.Sprintf("renderer all done.\n")))
	})
	return nil
}