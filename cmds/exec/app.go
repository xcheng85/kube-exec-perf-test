package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/xcheng85/kube-exec-perf-test/internal/config"
	"github.com/xcheng85/kube-exec-perf-test/internal/k8s/exec"
	"github.com/xcheng85/kube-exec-perf-test/internal/monolith"
	"github.com/xcheng85/kube-exec-perf-test/internal/worker"
	"golang.org/x/sync/errgroup"
	"net/http"
)

// composition root
// application in the hexongal arch
// app must implement monolith interface, which is required in each sub module
type app struct {
	config config.AppConfig
	// modules just like all the application controllers
	modules []monolith.Module
	mux     *chi.Mux
	exec    *exec.K8sExec
	// management multiple goroutines
	workerSyncer worker.WorkerSyncer
}

func (a *app) Config() config.AppConfig {
	return a.config
}

func (a *app) Mux() *chi.Mux {
	return a.mux
}

func (a *app) Exec() *exec.K8sExec {
	return a.exec
}

func (a *app) WorkerSyncer() worker.WorkerSyncer {
	return a.workerSyncer
}

// like the builder of ioc
func (a *app) startupModules() error {
	for _, module := range a.modules {
		if err := module.Startup(a.workerSyncer.Context(), a); err != nil {
			return err
		}
	}
	return nil
}

// worker for running Rest server for reverse proxy
func (a *app) runRest(ctx context.Context) error {
	// chi.Mux has the http.Handler embedded
	restServer := &http.Server{
		Addr:    a.config.Web.Address(),
		Handler: a.mux,
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Println("web server started")
		defer fmt.Println("web server shutdown")
		if err := restServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	group.Go(func() error {
		// received cancel signal from the derived
		<-gCtx.Done()
		fmt.Println("web server to be shutdown")
		// gracefully shut down rest server
		ctx, cancel := context.WithTimeout(context.Background(), a.config.ShutdownTimeout)
		defer cancel()
		if err := restServer.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	})
	// block here
	return group.Wait()
}
