package component

import "github.com/surkovvs/ct-app/appifaces"

type healthcheck Comp

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
