package component

import (
	"context"
	"fmt"

	"github.com/surkovvs/ct-app/appifaces"
)

type shutdown Comp

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

func (r shutdown) Shutdown(ctx context.Context) {
	defer r.prov.wgSd.Done()
	if err := r.object.(appifaces.Shutdowner).Shutdown(ctx); err != nil {
		r.prov.errChan <- fmt.Errorf(
			`group '%s', module '%s', shutdown: %w`,
			r.groupName, r.name, err)

		r.SetFailed()
		return
	}
	r.SetDone()
}
