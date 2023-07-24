package main

import (
	"github.com/sirupsen/logrus"
	"github.com/xcheng85/kube-exec-perf-test/internal/k8s/discovery"
	"github.com/xcheng85/kube-exec-perf-test/internal/k8s/exec"
	"golang.org/x/sync/errgroup"
	"time"
)

func createUnity(exec *exec.K8sExec) error{
	start := time.Now()
	//cmds := []string{"sh", "-c", "killall -s 0 entryPoint.sh"}
	cmds := []string{"sh", "-c", "ls"}
	stdout, stderr, err := exec.Exec(cmds)
	if err != nil {
		logrus.Errorf("Failed to exec:%v", err)
		return err
	}

	logrus.Infof("out:%s", stdout)
	logrus.Infof("err:%s", stderr)
	elapsed := time.Since(start)
	logrus.Infof("createUnity took %s", elapsed)
	return nil
}

func main() {
	gClientSet, gRestConfig, err := discovery.K8s()
	if err != nil {
		logrus.Fatalf("Could not get K8s.")
	}
	pod, container, namespace := "unity-test-app-evd-weu-demo1-qv3n4-pod", "rendering-engine", "evd-cia3dviz"
	k8s := exec.New(gClientSet, gRestConfig, pod, container, namespace)

	g := new(errgroup.Group)
	MAX_ITERATION := 100
	for i := 0; i < MAX_ITERATION; i++ {
		g.Go(func() error {
			err := createUnity(k8s)
			return err
		})
		// time.Sleep(time.Duration(5 * int(time.Second)))
	}
	if err := g.Wait(); err == nil {
		logrus.Println("All Done")
	}
}
