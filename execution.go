package ctapp

import (
	"context"
	"errors"
	"sync"

	"github.com/surkovvs/ct-app/component"
	"github.com/surkovvs/ct-app/compstor"
)

func (a *app) exec(ctx context.Context) {
	a.fillupRunnersWg()
	a.fillupShutdownWg()
	a.logger.Debug(`app started`, `application`, a.name)

	var initRunCtx context.Context
	initRunCtx, a.execution.initRunCancel = context.WithCancel(ctx)

	var initCtx context.Context
	if a.execution.initTimeout != nil {
		var cancel context.CancelFunc
		initCtx, cancel = context.WithTimeout(initRunCtx, *a.execution.initTimeout)
		defer cancel()
	} else {
		initCtx = initRunCtx
	}

	// backgroung group runs before others, all the components in background group runs concurrently
	group, err := a.storage.GetGroupByName(BackgroundGroup)
	if err != nil {
		if errors.Is(err, compstor.ErrGroupNotFound) {
			a.logger.Info(`background group not found`,
				"application", a.name)
		} else {
			a.logger.Error(`unexpected error`,
				"application", a.name,
				`group`, group.GetName(),
				"error", err)
		}
	} else {
		a.processBackground(initCtx, initRunCtx, group)
	}

	wg := sync.WaitGroup{}
	for _, group := range a.storage.GetOrderedGroupList() {
		if group.GetName() == BackgroundGroup {
			continue
		}
		wg.Add(1)
		go func(group compstor.SequentialGroup) {
			defer wg.Done()
			a.processInitializers(initCtx, group)
			a.processRunners(initRunCtx, group)
			a.processShutdowners(a.shutdown.ctx, group)
		}(group)
	}
	wg.Wait()
	a.execution.wg.Wait()

	close(a.execution.done)
}

func (a *app) fillupRunnersWg() {
	for _, group := range a.storage.GetOrderedGroupList() {
		for _, module := range group.GetComponents() {
			if module.IsRunner() || module.IsInitializer() {
				a.execution.wg.Add(1)
			}
		}
	}
}

func (a *app) processBackground(initCtx, runCtx context.Context, group compstor.SequentialGroup) {
	wg := sync.WaitGroup{}
	for _, module := range group.GetComponents() {
		wg.Add(1)
		go func(module component.Comp) {
			if module.Initializer().TrySetInProcess() {
				a.logger.Debug(`module initialization`,
					`application`, a.name,
					`group`, group.GetName(),
					`module`, module.Name())

				module.Initializer().Init(initCtx)
			}
			wg.Done()

			if (module.Initializer().IsDone() || !module.IsInitializer()) &&
				module.Runner().TrySetInProcess() {
				a.logger.Debug(`module running`,
					`application`, a.name,
					`group`, group.GetName(),
					`module`, module.Name())

				module.Runner().Run(runCtx)
			}

			if module.Runner().IsDone() && module.Shutdowner().TrySetInProcess() {
				a.logger.Debug(`module shutdown`,
					`application`, a.name,
					`group`, group.GetName(),
					`module`, module.Name())

				module.Shutdowner().Shutdown(a.shutdown.ctx)
			}
		}(module)
	}
	wg.Wait()
}

func (a *app) processInitializers(ctx context.Context, group compstor.SequentialGroup) {
	for _, module := range group.GetComponents() {
		if module.Initializer().TrySetInProcess() {

			a.logger.Debug(`module initialization`,
				`application`, a.name,
				`group`, group.GetName(),
				`module`, module.Name())

			module.Initializer().Init(ctx)
		}
	}
}

func (a *app) processRunners(ctx context.Context, group compstor.SequentialGroup) {
	for _, module := range group.GetComponents() {
		if (module.Initializer().IsDone() || !module.IsInitializer()) &&
			module.Runner().TrySetInProcess() {

			a.logger.Debug(`module running`,
				`application`, a.name,
				`group`, group.GetName(),
				`module`, module.Name())

			module.Runner().Run(ctx)
		}
	}
}

func (a *app) processShutdowners(ctx context.Context, group compstor.SequentialGroup) {
	for _, module := range group.GetComponents() {
		if module.Runner().IsDone() && module.Shutdowner().TrySetInProcess() {
			a.logger.Debug(`module shutdown`,
				`application`, a.name,
				`group`, group.GetName(),
				`module`, module.Name())
			module.Shutdowner().Shutdown(ctx)
		}
	}
}
