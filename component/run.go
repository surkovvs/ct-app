//nolint:dupl // TODO: to next refactoring
package component

import (
	"context"
	"fmt"

	"github.com/surkovvs/ct-app/appifaces"
)

type Run Comp

func (r Run) SetReady() {
	r.status.SetStatus(runReady, runMask)
}

func (r Run) SetInProcess() {
	r.status.SetStatus(runInProcess, runMask)
}

func (r Run) SetDone() {
	r.status.SetStatus(runDone, runMask)
}

func (r Run) SetFailed() {
	r.status.SetStatus(runFailed, runMask)
}

func (r Run) TrySetInProcess() bool {
	return r.status.TryChangeStatus(runReady, runInProcess, runMask)
}

func (r Run) IsReady() bool {
	return r.status.GetStatus().CompareMasked(runReady, runMask)
}

func (r Run) IsInProcess() bool {
	return r.status.GetStatus().CompareMasked(runInProcess, runMask)
}

func (r Run) IsDone() bool {
	return r.status.GetStatus().CompareMasked(runDone, runMask)
}

func (r Run) IsFailed() bool {
	return r.status.GetStatus().CompareMasked(runFailed, runMask)
}

func (r Run) Run(ctx context.Context) {
	defer r.prov.wgExec.Done()

	runner, ok := r.object.(appifaces.Runner)
	if !ok {
		panic(fmt.Sprintf(`group '%s', module '%s', incorrectly defined as Runner`, r.groupName, r.name))
	}
	if err := runner.Run(ctx); err != nil {
		r.prov.errChan <- fmt.Errorf(
			`group '%s', module '%s', run: %w`,
			r.groupName, r.name, err)

		r.SetFailed()
		return
	}
	r.SetDone()
}
