package component

import (
	"context"
	"fmt"

	"github.com/surkovvs/ct-app/appifaces"
)

type initialize Comp

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

func (r initialize) Init(ctx context.Context) {
	if err := r.object.(appifaces.Initializer).Init(ctx); err != nil {
		r.prov.errChan <- fmt.Errorf(
			`group '%s', module '%s', init: %w`,
			r.groupName, r.name, err)

		r.SetFailed()
		r.prov.wgExec.Done()
		return
	}
	if r.status.GetStatus().Querying(runMask) == 0 {
		r.prov.wgExec.Done()
	}
	r.SetDone()
}
