package component

import (
	"sync"

	"github.com/surkovvs/ct-app/appifaces"
	"github.com/surkovvs/ct-app/zorro"
)

type provider struct {
	wgExec  *sync.WaitGroup
	wgSd    *sync.WaitGroup
	errChan chan<- error
}
type hcDynamic struct {
	mu           *sync.Mutex
	hcProcessing chan struct{}
	hcErr        error
}

type Comp struct {
	prov      provider
	status    zorro.Zorro
	object    any
	hc        *hcDynamic
	groupName string
	name      string
}

const (
	initReady     zorro.Status = 1  // 0000000000000001
	initInProcess zorro.Status = 2  // 0000000000000010
	initDone      zorro.Status = 4  // 0000000000000100
	initFailed    zorro.Status = 8  // 0000000000001000
	initMask      zorro.Mask   = 15 // 0000000000001111

	runReady     zorro.Status = 16  // 0000000000010000
	runInProcess zorro.Status = 32  // 0000000000100000
	runDone      zorro.Status = 64  // 0000000001000000
	runFailed    zorro.Status = 128 // 0000000010000000
	runMask      zorro.Mask   = 240 // 0000000011110000

	shutdownReady     zorro.Status = 256  // 0000000100000000
	shutdownInProcess zorro.Status = 512  // 0000001000000000
	shutdownDone      zorro.Status = 1024 // 0000010000000000
	shutdownFailed    zorro.Status = 2048 // 0000100000000000
	shutdownMask      zorro.Mask   = 3840 // 0000111100000000

	healthcheckReady     zorro.Status = 4096  // 0001000000000000
	healthcheckInProcess zorro.Status = 8192  // 0010000000000000
	healthcheckDone      zorro.Status = 16384 // 0100000000000000
	healthcheckFailed    zorro.Status = 32768 // 1000000000000000
	healthcheckMask      zorro.Mask   = 61440 // 1111000000000000
)

type Define struct {
	GroupName  string
	CompName   string
	Component  any
	ExecWG     *sync.WaitGroup
	ShutdownWG *sync.WaitGroup
	ErrChan    chan<- error
}

func DefineComponent(d Define) Comp {
	status := zorro.New()
	if _, ok := d.Component.(appifaces.Healthchecker); ok {
		status.SetStatus(healthcheckReady, healthcheckMask)
	}
	if _, ok := d.Component.(appifaces.Initializer); ok {
		status.SetStatus(initReady, initMask)
	}
	if _, ok := d.Component.(appifaces.Runner); ok {
		status.SetStatus(runReady, runMask)
	}
	if _, ok := d.Component.(appifaces.Shutdowner); ok {
		status.SetStatus(shutdownReady, shutdownMask)
	}
	return Comp{
		name:   d.CompName,
		object: d.Component,
		status: status,
		prov: provider{
			wgExec:  d.ExecWG,
			wgSd:    d.ShutdownWG,
			errChan: d.ErrChan,
		},
		groupName: d.GroupName,
		hc:        &hcDynamic{mu: &sync.Mutex{}},
	}
}

func (c Comp) IsValid() bool {
	return c.status.GetStatus() != 0
}

func (c Comp) Name() string {
	return c.name
}

func (c Comp) IsHealthchecker() bool {
	return c.status.GetStatus().Querying(healthcheckMask) != 0
}

func (c Comp) Healthchecker() Healthcheck {
	return Healthcheck(c)
}

func (c Comp) IsInitializer() bool {
	return c.status.GetStatus().Querying(initMask) != 0
}

func (c Comp) Initializer() Initialize {
	return Initialize(c)
}

func (c Comp) IsRunner() bool {
	return c.status.GetStatus().Querying(runMask) != 0
}

func (c Comp) Runner() Run {
	return Run(c)
}

func (c Comp) IsShutdowner() bool {
	return c.status.GetStatus().Querying(shutdownMask) != 0
}

func (c Comp) Shutdowner() Shutdown {
	return Shutdown(c)
}
