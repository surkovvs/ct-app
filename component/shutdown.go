//nolint:dupl // TODO: to next refactoring
package component

import (
	"context"
	"fmt"

	"github.com/surkovvs/ct-app/appifaces"
)

type Shutdown Comp

func (r Shutdown) SetReady() {
	r.status.SetStatus(shutdownReady, shutdownMask)
}

func (r Shutdown) SetInProcess() {
	r.status.SetStatus(shutdownInProcess, shutdownMask)
}

func (r Shutdown) SetDone() {
	r.status.SetStatus(shutdownDone, shutdownMask)
}

func (r Shutdown) SetFailed() {
	r.status.SetStatus(shutdownFailed, shutdownMask)
}

func (r Shutdown) TrySetInProcess() bool {
	return r.status.TryChangeStatus(shutdownReady, shutdownInProcess, shutdownMask)
}

func (r Shutdown) IsReady() bool {
	return r.status.GetStatus().CompareMasked(shutdownReady, shutdownMask)
}

func (r Shutdown) IsInProcess() bool {
	return r.status.GetStatus().CompareMasked(shutdownInProcess, shutdownMask)
}

func (r Shutdown) IsDone() bool {
	return r.status.GetStatus().CompareMasked(shutdownDone, shutdownMask)
}

func (r Shutdown) IsFailed() bool {
	return r.status.GetStatus().CompareMasked(shutdownFailed, shutdownMask)
}

func (r Shutdown) Shutdown(ctx context.Context) {
	defer r.prov.wgSd.Done()

	shutdowner, ok := r.object.(appifaces.Shutdowner)
	if !ok {
		panic(fmt.Sprintf(`group '%s', module '%s', incorrectly defined as Shutdowner`, r.groupName, r.name))
	}
	if err := shutdowner.Shutdown(ctx); err != nil {
		r.prov.errChan <- fmt.Errorf(
			`group '%s', module '%s', shutdown: %w`,
			r.groupName, r.name, err)

		r.SetFailed()
		return
	}
	r.SetDone()
}
