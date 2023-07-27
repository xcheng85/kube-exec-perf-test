package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/xcheng85/kube-exec-perf-test/internal/k8s/exec"
	"golang.org/x/sync/errgroup"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

// Businiss logic regardless of api architecture
// DIP: ownership of interface
// define multiple custom type all at once
type (
	RendererService interface {
		Run(exec *exec.K8sExec, clientset kubernetes.Interface, ctx context.Context, e chan error)
	}
	// only expose interface
	rendererService struct {
	}
)

// do the interface checks
// if the AppImpl does not fulfill App interface, it will highlight
var _ RendererService = (*rendererService)(nil)

func NewRendererService() RendererService {
	return &rendererService{}
}

func (s rendererService) createUnity(exec *exec.K8sExec, ctx context.Context) error {
	select {
	case <-ctx.Done():
		logrus.Println("createUnity ctx done")
		return ctx.Err()
	default:
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
}

func (s rendererService) listPods(clientset kubernetes.Interface, ctx context.Context) error {
	select {
	case <-ctx.Done():
		logrus.Println("createUnity ctx done")
		return ctx.Err()
	default:
		start := time.Now()
		namespace := "evd-cia3dviz"
		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		logrus.Infof("There are %d pods in the cluster\n", len(pods.Items))
		elapsed := time.Since(start)
		logrus.Infof("listPods took %s", elapsed)
		return err
	}
}

func (s rendererService) runHeavyJob(exec *exec.K8sExec, clientset kubernetes.Interface, ctx context.Context) error {
	g := new(errgroup.Group)
	MAX_ITERATION := 100
	for i := 0; i < MAX_ITERATION; i++ {
		select {
		case <-ctx.Done():
			logrus.Println("runHeavyJob ctx done")
			return ctx.Err()
		default:
			g.Go(func() (err error) {
				if i > MAX_ITERATION {
					err = s.createUnity(exec, ctx)
				} else {
					err = s.listPods(clientset, ctx)
				}
				return err
			})
			time.Sleep(time.Duration(100 * int(time.Millisecond)))
		}
	}
	if err := g.Wait(); err == nil {
		logrus.Println("All Done")
	}
	return nil
}

func (s rendererService) Run(exec *exec.K8sExec, clientset kubernetes.Interface, ctx context.Context, e chan error) {
	err := s.runHeavyJob(exec, clientset, ctx)
	e <- err
	// select {
	// case <-ctx.Done():
	// 	logrus.Infof("ctx.Done()")
	// 	return
	// case <-s.runHeavyJob(exec, clientset):
	// 	return
	// }
}
