package ctapp

import (
	"time"

	"github.com/surkovvs/ct-app/vector"
)

func (a *App) gracefulShutdown() {
	defer close(a.shutdown.shutdownDone)

	go func() {
		time.Sleep(*a.shutdown.timeout)
		a.shutdown.ctxCancel()
	}()

	c := vector.NewConstructor(a.execution.reports)
	var vectors []any
	for _, module := range a.storage.GetUnsortedShutdowners() {
		vectors = append(vectors, module.ForceShutdown)
		// go func(module component.Comp) {
		// 	if rep := module.Shutdown(a.shutdown.ctx); rep.Code == component.CodeError {
		// 	}
		// 	if module.Shutdowner().TrySetInProcess() {
		// 		a.logger.Debug(`Module shutdown`,
		// 			`application`, a.name,
		// 			`module`, module.Name())

		// 		module.Shutdowner().ShutdownComponent(a.shutdown.ctx)
		// 	}
		// }(module)
	}
	c.Concurrently(a.shutdown.ctx, vectors...).Exec(a.shutdown.ctx)

	select {
	case <-a.shutdown.ctx.Done():
		a.reportUnfinished()
	default:
		a.logger.Info(`graceful shutdown finished`,
			"application", a.name)
	}
}

func (a *App) reportUnfinished() {
	var unfinishedIniters, unfinishedRunners, unfinishedShutdowners []string
	// for _, module := range a.storage.GetUnsortedInitializers() {
	// 	if module.Initializer().IsInProcess() {
	// 		unfinishedIniters = append(unfinishedIniters, module.Name())
	// 	}
	// }
	// for _, module := range a.storage.GetUnsortedRunners() {
	// 	if module.Runner().IsInProcess() {
	// 		unfinishedRunners = append(unfinishedRunners, module.Name())
	// 	}
	// }
	// for _, module := range a.storage.GetUnsortedShutdowners() {
	// 	if module.Shutdowner().IsInProcess() {
	// 		unfinishedShutdowners = append(unfinishedShutdowners, module.Name())
	// 	}
	// }

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
