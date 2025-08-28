package component

type TriggerError struct {
	w error
}

func (t TriggerError) Error() string {
	return t.w.Error()
}
