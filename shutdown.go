package ctapp

import (
	"time"

	"github.com/surkovvs/ct-app/component"
)

func (a *App) gracefulShutdown() {
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

				module.Shutdowner().ShutdownComponent(a.shutdown.ctx)
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
		a.reportUnfinished()
	}
}

func (a *App) reportUnfinished() {
	var unfinishedIniters, unfinishedRunners, unfinishedShutdowners []string
	for _, module := range a.storage.GetUnsortedInitializers() {
		if module.Initializer().IsInProcess() {
			unfinishedIniters = append(unfinishedIniters, module.Name())
		}
	}
	for _, module := range a.storage.GetUnsortedRunners() {
		if module.Runner().IsInProcess() {
			unfinishedRunners = append(unfinishedRunners, module.Name())
		}
	}
	for _, module := range a.storage.GetUnsortedShutdowners() {
		if module.Shutdowner().IsInProcess() {
			unfinishedShutdowners = append(unfinishedShutdowners, module.Name())
		}
	}

	a.logger.Error(`graceful shutdown timeout exeeded, got unfinished modules`,
		"application", a.name,
		"initialization", unfinishedIniters,
		"running", unfinishedRunners,
		"shutdown", unfinishedShutdowners,
	)

	// switch {
	// case len(unfinishedIniters) != 0 && len(unfinishedRunners) != 0 && len(unfinishedShutdowners) != 0:
	// 	a.logger.Error(`graceful shutdown timeout exeeded, got unfinished modules`,
	// 		"application", a.name,
	// 		"initialization", unfinishedIniters,
	// 		"running", unfinishedRunners,
	// 		"shutdown", unfinishedShutdowners,
	// 	)
	// case len(unfinishedIniters) != 0 && len(unfinishedRunners) != 0:
	// 	a.logger.Error(`graceful shutdown timeout exeeded, unfinished modules`,
	// 		"application", a.name,
	// 		"initialization", unfinishedIniters,
	// 		"running", unfinishedRunners,
	// 	)
	// case len(unfinishedIniters) != 0 && len(unfinishedShutdowners) != 0:
	// 	a.logger.Error(`graceful shutdown timeout exeeded, got unfinished modules`,
	// 		"application", a.name,
	// 		"initialization", unfinishedIniters,
	// 		"shutdown", unfinishedShutdowners,
	// 	)
	// case len(unfinishedRunners) != 0 && len(unfinishedShutdowners) != 0:
	// 	a.logger.Error(`graceful shutdown timeout exeeded, got unfinished modules`,
	// 		"application", a.name,
	// 		"running", unfinishedRunners,
	// 		"shutdown", unfinishedShutdowners,
	// 	)
	// case len(unfinishedIniters) != 0:
	// 	a.logger.Error(`graceful shutdown timeout exeeded, got unfinished modules`,
	// 		"application", a.name,
	// 		"initialization", unfinishedIniters,
	// 	)
	// case len(unfinishedRunners) != 0:
	// 	a.logger.Error(`graceful shutdown timeout exeeded, got unfinished modules`,
	// 		"application", a.name,
	// 		"running", unfinishedRunners,
	// 	)
	// case len(unfinishedShutdowners) != 0:
	// 	a.logger.Error(`graceful shutdown timeout exeeded, got unfinished modules`,
	// 		"application", a.name,
	// 		"shutdown", unfinishedShutdowners,
	// 	)
	// default:
	// 	a.logger.Error(`graceful shutdown timeout exeeded`,
	// 		"application", a.name,
	// 	)
	// }
}
