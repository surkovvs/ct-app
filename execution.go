package ctapp

import (
	"context"
	"errors"
	"sync"

	"github.com/surkovvs/ct-app/component"
	"github.com/surkovvs/ct-app/compstor"
)

func (a *App) fillupRunnersWg() {
	for _, group := range a.storage.GetOrderedGroupList() {
		for _, module := range group.GetComponents() {
			if module.IsRunner() || module.IsInitializer() {
				a.execution.wg.Add(1)
			}
		}
	}
}

func (a *App) fillupShutdownWg() {
	for range a.storage.GetUnsortedShutdowners() {
		a.shutdown.wg.Add(1)
	}
}

func (a *App) exec() {
	// backgroung groups runs before others, all the components in background group runs concurrently
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.processBackground()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.processBackgroundSync()
	}()
	wg.Wait()

	a.processSequentialGroups()

	a.execution.wg.Wait()
	close(a.execution.done)
}

func (a *App) processBackground() {
	bgGroup, err := a.storage.GetGroupByName(BackgroundGroup)
	if err != nil {
		if errors.Is(err, compstor.ErrGroupNotFound) {
			a.logger.Debug(`background group not found`,
				"application", a.name)
		} else {
			a.logger.Error(`unexpected error`,
				"application", a.name,
				`group`, bgGroup.GetName(),
				"error", err)
		}
	} else {
		wg := sync.WaitGroup{}
		for _, module := range bgGroup.GetComponents() {
			wg.Add(1)
			go func(module component.Comp) {
				a.processInitializer(a.execution.initCtx, module)
				wg.Done()
				a.processRunner(a.execution.runCtx, module)
				a.processShutdowner(a.shutdown.ctx, module)
			}(module)
		}
		wg.Wait()
	}
}

func (a *App) processBackgroundSync() {
	bgsGroup, err := a.storage.GetGroupByName(BackgroundSyncGroup)
	if err != nil {
		if errors.Is(err, compstor.ErrGroupNotFound) {
			a.logger.Debug(`background sync group not found`,
				"application", a.name)
		} else {
			a.logger.Error(`unexpected error`,
				"application", a.name,
				`group`, bgsGroup.GetName(),
				"error", err)
		}
	} else {
		wg := sync.WaitGroup{}
		for _, module := range bgsGroup.GetComponents() {
			wg.Add(1)
			go func(module component.Comp) {
				a.processInitializer(a.execution.initCtx, module)
				wg.Done()
				wg.Wait()
				a.processRunner(a.execution.runCtx, module)
				a.processShutdowner(a.shutdown.ctx, module)
			}(module)
		}
		wg.Wait()
	}
}

func (a *App) processSequentialGroups() {
	wg := sync.WaitGroup{}
	for _, group := range a.storage.GetOrderedGroupList() {
		if group.GetName() == BackgroundGroup {
			continue
		}
		wg.Add(1)
		go func(group compstor.SequentialGroup) {
			defer wg.Done()
			for _, module := range group.GetComponents() {
				a.processInitializer(a.execution.initCtx, module)
			}
			for _, module := range group.GetComponents() {
				a.processRunner(a.execution.runCtx, module)
			}
			for _, module := range group.GetComponents() {
				a.processShutdowner(a.shutdown.ctx, module)
			}
		}(group)
	}
	wg.Wait()
}

func (a *App) processInitializer(ctx context.Context, module component.Comp) {
	if module.Initializer().TrySetInProcess() && ctx.Err() == nil {
		a.logger.Debug(`module initialization`,
			`application`, a.name,
			`group`, module.GroupName(),
			`module`, module.Name())

		module.Initializer().InitComponent(ctx)
	}
}

func (a *App) processRunner(ctx context.Context, module component.Comp) {
	if (module.Initializer().IsDone() || !module.IsInitializer()) &&
		module.Runner().TrySetInProcess() && ctx.Err() == nil {
		a.logger.Debug(`module running`,
			`application`, a.name,
			`group`, module.GroupName(),
			`module`, module.Name())

		module.Runner().RunComponent(ctx)
	} else if ctx.Err() != nil && (module.IsRunner() && module.Initializer().IsDone()) {
		a.execution.wg.Done()
	}
}

func (a *App) processShutdowner(ctx context.Context, module component.Comp) {
	if module.Runner().IsDone() && module.Shutdowner().TrySetInProcess() && ctx.Err() == nil {
		a.logger.Debug(`module shutdown`,
			`application`, a.name,
			`group`, module.GroupName(),
			`module`, module.Name())

		module.Shutdowner().ShutdownComponent(ctx)
	}
}
