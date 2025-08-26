package component

import (
	"context"
	"fmt"

	"github.com/surkovvs/ct-app/appifaces"
)

type run Comp

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

func (r run) Run(ctx context.Context) {
	defer r.prov.wgExec.Done()
	if err := r.object.(appifaces.Runner).Run(ctx); err != nil {
		r.prov.errChan <- fmt.Errorf(
			`group '%s', module '%s', run: %w`,
			r.groupName, r.name, err)

		r.SetFailed()
		return
	}
	r.SetDone()
}
