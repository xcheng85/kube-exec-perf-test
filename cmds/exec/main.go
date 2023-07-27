package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/xcheng85/kube-exec-perf-test/internal/config"
	"github.com/xcheng85/kube-exec-perf-test/internal/monolith"
	"github.com/xcheng85/kube-exec-perf-test/internal/worker"
	"github.com/xcheng85/kube-exec-perf-test/renderer"
	"os"

	"github.com/xcheng85/kube-exec-perf-test/internal/k8s/discovery"
	"github.com/xcheng85/kube-exec-perf-test/internal/k8s/exec"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run() (err error) {
	// set up driven adapters
	// create deps
	config, err := config.NewAppConfig()
	if err != nil {
		return err
	}
	mux := createMux()

	gClientSet, gRestConfig, err := discovery.K8s()
	if err != nil {
		logrus.Fatalf("Could not get K8s.")
	}
	pod, container, namespace := "unity-test-app-evd-weu-demo1-bis2v-pod", "rendering-engine", "evd-cia3dviz"
	exec := exec.New(gClientSet, gRestConfig, pod, container, namespace)

	workerSyncer := worker.NewSyncer()
	modules := []monolith.Module{
		&players.Module{},
	}
	// setup application
	// build the app with deps
	myapp := app{
		config:       config,
		modules:      modules,
		mux:          mux,
		exec:         exec,
		workerSyncer: workerSyncer,
	}
	// set up Driver adapters
	// bind rest and grpc routes
	if err = myapp.startupModules(); err != nil {
		return err
	}
	fmt.Println("started tennis-shop application")
	defer fmt.Println("stopped tennis-shop application")

	// blocking main thread
	myapp.workerSyncer.Add(
		myapp.runRest,
	)
	return myapp.workerSyncer.Sync()
}

func createMux() *chi.Mux {
	return chi.NewMux()
}
