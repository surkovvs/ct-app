package component

import (
	"context"
	"fmt"

	"github.com/surkovvs/ct-app/appifaces"
)

type Healthcheck Comp

func (r Healthcheck) SetReady() {
	r.status.SetStatus(healthcheckReady, healthcheckMask)
}

func (r Healthcheck) SetInProcess() {
	r.status.SetStatus(healthcheckInProcess, healthcheckMask)
}

func (r Healthcheck) SetDone() {
	r.status.SetStatus(healthcheckDone, healthcheckMask)
}

func (r Healthcheck) SetFailed() {
	r.status.SetStatus(healthcheckFailed, healthcheckMask)
}

func (r Healthcheck) TrySetInProcess() bool {
	switch {
	case r.status.TryChangeStatus(healthcheckReady, healthcheckInProcess, healthcheckMask):
		return true
	case r.status.TryChangeStatus(healthcheckDone, healthcheckInProcess, healthcheckMask):
		return true
	case r.status.TryChangeStatus(healthcheckFailed, healthcheckInProcess, healthcheckMask):
		return true
	}
	return false
}

func (r Healthcheck) IsReady() bool {
	return r.status.GetStatus().CompareMasked(healthcheckReady, healthcheckMask)
}

func (r Healthcheck) IsInProcess() bool {
	return r.status.GetStatus().CompareMasked(healthcheckInProcess, healthcheckMask)
}

func (r Healthcheck) IsDone() bool {
	return r.status.GetStatus().CompareMasked(healthcheckDone, healthcheckMask)
}

func (r Healthcheck) IsFailed() bool {
	return r.status.GetStatus().CompareMasked(healthcheckFailed, healthcheckMask)
}

func (r Healthcheck) Healthcheck(ctx context.Context) error {
	if !r.IsInProcess() {
		r.hc.mu.Lock()
		defer r.hc.mu.Unlock()

		r.hc.hcProcessing = make(chan struct{})
		defer close(r.hc.hcProcessing)

		r.SetInProcess()
		healthchecker, ok := r.object.(appifaces.Healthchecker)
		if !ok {
			panic(fmt.Sprintf(`group '%s', module '%s', incorrectly defined as Healthchecker`, r.groupName, r.name))
		}
		if err := healthchecker.Healthcheck(ctx); err != nil {
			r.prov.errChan <- fmt.Errorf(
				`group '%s', module '%s', healthcheck: %w`,
				r.groupName, r.name, err)

			r.hc.hcErr = err
			r.SetFailed()
			return err
		}

		r.SetDone()
		return nil
	}

	<-r.hc.hcProcessing
	return r.hc.hcErr
}
