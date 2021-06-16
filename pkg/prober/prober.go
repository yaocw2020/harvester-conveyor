package prober

import (
	"context"
	"time"

	"github.com/tevino/tcp-shaker"
	"k8s.io/klog/v2"
)

type Prober interface {
	Probe(address string, timeout time.Duration) error
}

type tcpProber struct {
	*tcp.Checker
}

func (t *tcpProber) Probe(address string, timeout time.Duration) error {
	return t.CheckAddr(address, timeout)
}

func newTCPProber(ctx context.Context) *tcpProber {
	checker := tcp.NewChecker()
	go func() {
		if err := checker.CheckingLoop(ctx); err != nil {
			klog.Errorf("checking loop stopped due to fatal error: %s", err.Error())
		}
	}()

	<-checker.WaitReady()

	return &tcpProber{Checker: checker}
}
