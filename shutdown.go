package ctapp

import (
	"time"

	"github.com/surkovvs/ct-app/component"
)

func (a *app) fillupShutdownWg() {
	for range a.storage.GetUnsortedShutdowners() {
		a.shutdown.wg.Add(1)
	}
}

func (a *app) gracefulShutdown() {
	defer close(a.shutdown.shutdownDone)

	go func() {
		time.Sleep(*a.shutdown.timeout)
		a.shutdown.ctxCancel()
	}()

	for _, module := range a.storage.GetUnsortedShutdowners() {
		go func(module component.Comp) {
			if module.Shutdowner().TrySetInProcess() {
				a.logger.Debug(`Module shutdown`,
					`application`, a.name,
					`module`, module.Name())

				module.Shutdowner().Shutdown(a.shutdown.ctx)
			}
		}(module)
	}

	gsDone := make(chan struct{})
	go func() {
		a.shutdown.wg.Wait()
		close(gsDone)
	}()

	select {
	case <-gsDone:
		a.logger.Info(`graceful shutdown finished`,
			"application", a.name)

	case <-a.shutdown.ctx.Done():
		var unfinished []string
		for _, module := range a.storage.GetUnsortedShutdowners() {
			if module.Shutdowner().IsInProcess() {
				unfinished = append(unfinished, module.Name())
			}
		}
		if unfinished != nil {
			a.logger.Error(`graceful shutdown timeout exeeded`,
				"application", a.name,
				"unfinished shutdowns for modules", unfinished,
			)
		}
	}
}
