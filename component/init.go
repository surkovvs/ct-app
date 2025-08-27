package component

import (
	"context"
	"fmt"

	"github.com/surkovvs/ct-app/appifaces"
)

type Initialize Comp

func (r Initialize) SetReady() {
	r.status.SetStatus(initReady, initMask)
}

func (r Initialize) SetInProcess() {
	r.status.SetStatus(initInProcess, initMask)
}

func (r Initialize) SetDone() {
	r.status.SetStatus(initDone, initMask)
}

func (r Initialize) SetFailed() {
	r.status.SetStatus(initFailed, initMask)
}

func (r Initialize) TrySetInProcess() bool {
	return r.status.TryChangeStatus(initReady, initInProcess, initMask)
}

func (r Initialize) IsReady() bool {
	return r.status.GetStatus().CompareMasked(initReady, initMask)
}

func (r Initialize) IsInProcess() bool {
	return r.status.GetStatus().CompareMasked(initInProcess, initMask)
}

func (r Initialize) IsDone() bool {
	return r.status.GetStatus().CompareMasked(initDone, initMask)
}

func (r Initialize) IsFailed() bool {
	return r.status.GetStatus().CompareMasked(initFailed, initMask)
}

func (r Initialize) Init(ctx context.Context) {
	initializer, ok := r.object.(appifaces.Initializer)
	if !ok {
		panic(fmt.Sprintf(`group '%s', module '%s', incorrectly defined as Initializer`, r.groupName, r.name))
	}
	if err := initializer.Init(ctx); err != nil {
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
