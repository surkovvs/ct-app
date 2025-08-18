package component

import (
	"github.com/surkovvs/ct-app/appifaces"
	"github.com/surkovvs/ct-app/zorro"
)

type Comp struct {
	name   string
	object any
	status zorro.Zorro
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

type (
	healthcheck Comp
	initialize  Comp
	run         Comp
	shutdown    Comp
)

func DefineComponent(name string, component any) Comp {
	status := zorro.New()
	if _, ok := component.(appifaces.Healthchecker); ok {
		status.SetStatus(healthcheckReady, healthcheckMask)
	}
	if _, ok := component.(appifaces.Initializer); ok {
		status.SetStatus(initReady, initMask)
	}
	if _, ok := component.(appifaces.Runner); ok {
		status.SetStatus(runReady, runMask)
	}
	if _, ok := component.(appifaces.Shutdowner); ok {
		status.SetStatus(shutdownReady, shutdownMask)
	}
	return Comp{
		name:   name,
		object: component,
		status: status,
	}
}

func (c Comp) IsValid() bool {
	return c.status.GetStatus() != 0
}

func (c Comp) Name() string {
	return c.name
}

// healthcheck crew

func (c Comp) IsHealthchecker() bool {
	return c.status.GetStatus().Querying(healthcheckMask) != 0
}

func (c Comp) Healthchecker() healthcheck {
	return healthcheck(c)
}

func (r healthcheck) SetReady() {
	r.status.SetStatus(healthcheckReady, healthcheckMask)
}

func (r healthcheck) SetInProcess() {
	r.status.SetStatus(healthcheckInProcess, healthcheckMask)
}

func (r healthcheck) SetDone() {
	r.status.SetStatus(healthcheckDone, healthcheckMask)
}

func (r healthcheck) SetFailed() {
	r.status.SetStatus(healthcheckFailed, healthcheckMask)
}

func (r healthcheck) TrySetInProcess() bool {
	return r.status.TryChangeStatus(healthcheckReady, healthcheckInProcess, healthcheckMask)
}

func (r healthcheck) IsReady() bool {
	return r.status.GetStatus().CompareMasked(healthcheckReady, healthcheckMask)
}

func (r healthcheck) IsInProcess() bool {
	return r.status.GetStatus().CompareMasked(healthcheckInProcess, healthcheckMask)
}

func (r healthcheck) IsDone() bool {
	return r.status.GetStatus().CompareMasked(healthcheckDone, healthcheckMask)
}

func (r healthcheck) IsFailed() bool {
	return r.status.GetStatus().CompareMasked(healthcheckFailed, healthcheckMask)
}

func (r healthcheck) Get() appifaces.Healthchecker {
	return r.object.(appifaces.Healthchecker)
}

// initialize crew

func (c Comp) IsInitializer() bool {
	return c.status.GetStatus().Querying(initMask) != 0
}

func (c Comp) Initializer() initialize {
	return initialize(c)
}

func (r initialize) SetReady() {
	r.status.SetStatus(initReady, initMask)
}

func (r initialize) SetInProcess() {
	r.status.SetStatus(initInProcess, initMask)
}

func (r initialize) SetDone() {
	r.status.SetStatus(initDone, initMask)
}

func (r initialize) SetFailed() {
	r.status.SetStatus(initFailed, initMask)
}

func (r initialize) TrySetInProcess() bool {
	return r.status.TryChangeStatus(initReady, initInProcess, initMask)
}

func (r initialize) IsReady() bool {
	return r.status.GetStatus().CompareMasked(initReady, initMask)
}

func (r initialize) IsInProcess() bool {
	return r.status.GetStatus().CompareMasked(initInProcess, initMask)
}

func (r initialize) IsDone() bool {
	return r.status.GetStatus().CompareMasked(initDone, initMask)
}

func (r initialize) IsFailed() bool {
	return r.status.GetStatus().CompareMasked(initFailed, initMask)
}

func (r initialize) Get() appifaces.Initializer {
	return r.object.(appifaces.Initializer)
}

// run crew

func (c Comp) IsRunner() bool {
	return c.status.GetStatus().Querying(runMask) != 0
}

func (c Comp) Runner() run {
	return run(c)
}

func (r run) SetReady() {
	r.status.SetStatus(runReady, runMask)
}

func (r run) SetInProcess() {
	r.status.SetStatus(runInProcess, runMask)
}

func (r run) SetDone() {
	r.status.SetStatus(runDone, runMask)
}

func (r run) SetFailed() {
	r.status.SetStatus(runFailed, runMask)
}

func (r run) TrySetInProcess() bool {
	return r.status.TryChangeStatus(runReady, runInProcess, runMask)
}

func (r run) IsReady() bool {
	return r.status.GetStatus().CompareMasked(runReady, runMask)
}

func (r run) IsInProcess() bool {
	return r.status.GetStatus().CompareMasked(runInProcess, runMask)
}

func (r run) IsDone() bool {
	return r.status.GetStatus().CompareMasked(runDone, runMask)
}

func (r run) IsFailed() bool {
	return r.status.GetStatus().CompareMasked(runFailed, runMask)
}

func (r run) Get() appifaces.Runner {
	return r.object.(appifaces.Runner)
}

// shutdown crew

func (c Comp) IsShutdowner() bool {
	return c.status.GetStatus().Querying(shutdownMask) != 0
}

func (c Comp) Shutdowner() shutdown {
	return shutdown(c)
}

func (r shutdown) SetReady() {
	r.status.SetStatus(shutdownReady, shutdownMask)
}

func (r shutdown) SetInProcess() {
	r.status.SetStatus(shutdownInProcess, shutdownMask)
}

func (r shutdown) SetDone() {
	r.status.SetStatus(shutdownDone, shutdownMask)
}

func (r shutdown) SetFailed() {
	r.status.SetStatus(shutdownFailed, shutdownMask)
}

func (r shutdown) TrySetInProcess() bool {
	return r.status.TryChangeStatus(shutdownReady, shutdownInProcess, shutdownMask)
}

func (r shutdown) IsReady() bool {
	return r.status.GetStatus().CompareMasked(shutdownReady, shutdownMask)
}

func (r shutdown) IsInProcess() bool {
	return r.status.GetStatus().CompareMasked(shutdownInProcess, shutdownMask)
}

func (r shutdown) IsDone() bool {
	return r.status.GetStatus().CompareMasked(shutdownDone, shutdownMask)
}

func (r shutdown) IsFailed() bool {
	return r.status.GetStatus().CompareMasked(shutdownFailed, shutdownMask)
}

func (r shutdown) Get() appifaces.Shutdowner {
	return r.object.(appifaces.Shutdowner)
}
