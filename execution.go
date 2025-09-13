package ctapp

import (
	"context"

	"github.com/surkovvs/ct-app/component"
	"github.com/surkovvs/ct-app/vector"
)

// func (a *App) fillupRunnersWg() {
// 	for _, group := range a.storage.GetOrderedGroupList() {
// 		for _, module := range group.GetComponents() {
// 			if module.IsRunner() || module.IsInitializer() {
// 				a.execution.wg.Add(1)
// 			}
// 		}
// 	}
// }

// func (a *App) fillupShutdownWg() {
// 	for range a.storage.GetUnsortedShutdowners() {
// 		a.shutdown.wg.Add(1)
// 	}
// }

func (a *App) exec() {
	// backgroung groups runs before others, all the components in background group runs concurrently
	// wg := sync.WaitGroup{}
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	a.processBackground()
	// }()
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	a.processBackgroundSync()
	// }()
	// wg.Wait()

	a.processSequentialGroups()

	a.execution.wg.Wait()
	close(a.execution.done)
}

func (a *App) processBackground() {
	// bgGroup, err := a.storage.GetGroupByName(BackgroundGroup)
	// if err != nil {
	// 	if errors.Is(err, compstor.ErrGroupNotFound) {
	// 		a.logger.Debug(`background group not found`,
	// 			"application", a.name)
	// 	} else {
	// 		a.logger.Error(`unexpected error`,
	// 			"application", a.name,
	// 			`group`, bgGroup.GetName(),
	// 			"error", err)
	// 	}
	// } else {
	// 	wg := sync.WaitGroup{}
	// 	for _, module := range bgGroup.GetComponents() {
	// 		wg.Add(1)
	// 		go func(module component.Comp) {
	// 			a.processInitializer(a.execution.initCtx, module)
	// 			wg.Done()
	// 			a.processRunner(a.execution.runCtx, module)
	// 			a.processShutdowner(a.shutdown.ctx, module)
	// 		}(module)
	// 	}
	// 	wg.Wait()
	// }
}

func (a *App) processBackgroundSync() {
	// bgsGroup, err := a.storage.GetGroupByName(BackgroundSyncGroup)
	// if err != nil {
	// 	if errors.Is(err, compstor.ErrGroupNotFound) {
	// 		a.logger.Debug(`background sync group not found`,
	// 			"application", a.name)
	// 	} else {
	// 		a.logger.Error(`unexpected error`,
	// 			"application", a.name,
	// 			`group`, bgsGroup.GetName(),
	// 			"error", err)
	// 	}
	// } else {
	// 	wg := sync.WaitGroup{}
	// 	for _, module := range bgsGroup.GetComponents() {
	// 		c := vector.NewConstructor(a.execution.reports)
	// 		c.Sequentially()
	// 		Sequentially(a.execution.runCtx, module.Init().InitComponent(a.execution.initCtx))
	// 		wg.Add(1)
	// 		go func(module component.Comp) {
	// 			a.processInitializer(a.execution.initCtx, module)
	// 			wg.Done()
	// 			wg.Wait()
	// 			a.processRunner(a.execution.runCtx, module)
	// 			a.processShutdowner(a.shutdown.ctx, module)
	// 		}(module)
	// 	}
	// 	wg.Wait()
	// }
}

func (a *App) processSequentialGroups() {
	c := vector.NewConstructor(a.execution.reports)
	var vectors []any
	for _, group := range a.storage.GetOrderedGroupList() {
		if group.GetName() == BackgroundGroup || group.GetName() == BackgroundSyncGroup {
			continue
		}

		var inits, runs, sds []any
		for _, module := range group.GetComponents() {
			inits = append(inits, a.callback(module, "initialization"), module.Init)
			runs = append(runs, a.callback(module, "running"), module.Run)
			sds = append(sds, a.callback(module, "shutdown"), module.Shutdown)
		}
		vectors = append(vectors, c.Sequentially(a.execution.runCtx,
			c.Sequentially(a.execution.initCtx, inits...),
			c.Sequentially(a.execution.runCtx, runs...),
			c.Sequentially(a.shutdown.ctx, sds...),
		))
	}
	c.Concurrently(a.execution.runCtx, vectors...).Exec(a.execution.runCtx)
}

func (a *App) callback(module component.Comp, stage string) func(ctx context.Context) component.Report {
	return func(ctx context.Context) component.Report {
		a.logger.Debug(`module `+stage,
			`application`, a.name,
			`group`, module.GroupName(),
			`module`, module.Name())
		return component.Report{}
	}
}
