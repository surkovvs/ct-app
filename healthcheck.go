package ctapp

import (
	"context"
	"sync"

	"github.com/surkovvs/ct-app/component"
)

func (a *App) Healthcheck(ctx context.Context) []error {
	var errs []error
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	for _, module := range a.storage.GetUnsortedHealthcheckers() {
		wg.Add(1)
		go func(module component.Comp) {
			defer wg.Done()
			if err := module.Healthchecker().Healthcheck(ctx); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}(module)
	}
	wg.Wait()

	return errs
}
