package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/xcheng85/kube-exec-perf-test/internal/k8s/exec"
	"golang.org/x/sync/errgroup"
	"time"
	// "github.com/xcheng85/Go-EDA/players/internal/domain"
	// "github.com/xcheng85/Go-EDA/players/internal/dto"
)

// Businiss logic regardless of api architecture
// DIP: ownership of interface
// define multiple custom type all at once
type (
	RendererService interface {
		Run(exec *exec.K8sExec, ctx context.Context) error
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

func (s rendererService) createUnity(exec *exec.K8sExec) error {
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

func (s rendererService) Run(exec *exec.K8sExec, ctx context.Context) error {
	g := new(errgroup.Group)
	MAX_ITERATION := 30
	for i := 0; i < MAX_ITERATION; i++ {
		g.Go(func() (err error) {
			err = s.createUnity(exec)
			return err
		})
		// time.Sleep(time.Duration(5 * int(time.Second)))
	}
	if err := g.Wait(); err == nil {
		logrus.Println("All Done")
	}
	return nil
}
