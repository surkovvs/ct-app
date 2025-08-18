package ctapp

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/surkovvs/ct-app/component"
)

func (a *app) gracefulShutdown() {
	wg := sync.WaitGroup{}

	gsDone := make(chan struct{})

	ctx, cancel := context.WithTimeout(a.shutdown.ctx, *a.shutdown.timeout)
	defer cancel()

	for _, module := range a.storage.GetUnsortedShutdowners() {
		wg.Add(1)
		go func(module component.Comp) {
			defer wg.Done()

			if module.Shutdowner().TrySetInProcess() {
				if err := module.Shutdowner().Get().Shutdown(ctx); err != nil {
					a.execution.errFlow <- fmt.Errorf(
						`shutdown module "%s", failed: %w`,
						module.Name(), err)
					module.Shutdowner().SetFailed()
					return
				}
				module.Shutdowner().SetDone()
			}
		}(module)
	}

	go func() {
		wg.Wait()
		close(gsDone)
	}()

SelLabel:
	select {
	case <-gsDone:
		a.logger.Info(`graceful shutdown finished for all bacground runners`,
			"application", a.name)
		for _, module := range a.storage.GetUnsortedShutdowners() {
			if module.Shutdowner().IsInProcess() {
				a.logger.Info(`graceful shutdown for cron module still in process, app shutdown time goes to max limit`,
					"application", a.name,
					`module`, module.Name())
				<-ctx.Done()
				break SelLabel
			}
		}

	case <-ctx.Done():
		a.logger.Error(`graceful shutdown timeout exeeded`,
			"application", a.name)
	}

	os.Exit(a.shutdown.exitCode)
}
