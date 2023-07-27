package players

import (
	"context"
	"github.com/xcheng85/kube-exec-perf-test/internal/monolith"
	"github.com/xcheng85/kube-exec-perf-test/renderer/internal/rest"
	"github.com/xcheng85/kube-exec-perf-test/renderer/internal/service"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	rendererService := service.NewRendererService()
	if err := rest.RegisterServer(rendererService, ctx, mono.Mux(), mono.Exec()); err != nil {
		return err
	}
	return nil
}