package component

import "github.com/surkovvs/ct-app/zorro"

const (
	ready     zorro.Status = 4369  // 0001000100010001
	inProcess zorro.Status = 8738  // 0010001000100010
	done      zorro.Status = 17476 // 0100010001000100
	failed    zorro.Status = 34952 // 1000100010001000

	// initReady     zorro.Status = 1  // 0000000000000001
	// initInProcess zorro.Status = 2  // 0000000000000010
	// initDone      zorro.Status = 4  // 0000000000000100
	// initFailed    zorro.Status = 8  // 0000000000001000
	initMask zorro.Mask = 15 // 0000000000001111

	// runReady     zorro.Status = 16  // 0000000000010000
	// runInProcess zorro.Status = 32  // 0000000000100000
	// runDone      zorro.Status = 64  // 0000000001000000
	// runFailed    zorro.Status = 128 // 0000000010000000
	runMask zorro.Mask = 240 // 0000000011110000

	// shutdownReady     zorro.Status = 256  // 0000000100000000
	// shutdownInProcess zorro.Status = 512  // 0000001000000000
	// shutdownDone      zorro.Status = 1024 // 0000010000000000
	// shutdownFailed    zorro.Status = 2048 // 0000100000000000
	shutdownMask zorro.Mask = 3840 // 0000111100000000

	// healthcheckReady     zorro.Status = 4096  // 0001000000000000
	// healthcheckInProcess zorro.Status = 8192  // 0010000000000000
	// healthcheckDone      zorro.Status = 16384 // 0100000000000000
	// healthcheckFailed    zorro.Status = 32768 // 1000000000000000
	healthcheckMask zorro.Mask = 61440 // 1111000000000000
)

var namedStatuses = map[uint64]string{
	1: "ready",
	2: "in_process",
	4: "done",
	8: "failed",
}

type statusProvider struct {
	provided zorro.Mask
	comp     Comp
}

func (r statusProvider) setReady() {
	r.comp.status.SetStatus(ready, r.provided)
}

func (r statusProvider) setInProcess() {
	r.comp.status.SetStatus(inProcess, r.provided)
}

func (r statusProvider) setDone() {
	r.comp.status.SetStatus(done, r.provided)
}

func (r statusProvider) setFailed() {
	r.comp.status.SetStatus(failed, r.provided)
}

func (r statusProvider) trySetInProcess() bool {
	return r.comp.status.TryChangeStatus(ready, inProcess, r.provided)
}

func (r statusProvider) isReady() bool {
	return r.comp.status.GetStatus().CompareMasked(ready, r.provided)
}

func (r statusProvider) isInProcess() bool {
	return r.comp.status.GetStatus().CompareMasked(inProcess, r.provided)
}

func (r statusProvider) isDone() bool {
	return r.comp.status.GetStatus().CompareMasked(done, r.provided)
}

func (r statusProvider) isFailed() bool {
	return r.comp.status.GetStatus().CompareMasked(failed, r.provided)
}

func (r statusProvider) namedStatus() string {
	bald := zorro.Status(r.comp.status.GetStatus().Querying(r.provided))
	s, ok := namedStatuses[bald.ShiftTrailingZeros(r.provided)]
	if !ok {
		return "unknown"
	}
	return s
}
