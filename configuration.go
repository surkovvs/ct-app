package ctapp

import (
	"os"
	"time"

	"github.com/surkovvs/ct-app/appifaces"
)

type ConfigApp struct {
	Silient                  bool
	KeepProcessOnModuleError bool // TODO: добавить реализацию
	Name                     *string
	InitTimeout              *time.Duration
	ShutdownTimeout          *time.Duration
}

func (c ConfigApp) IsSilientMode() bool {
	return c.Silient
}

func (c ConfigApp) GetApplicationName() *string {
	return c.Name
}

func (c ConfigApp) GetInitTimeout() *time.Duration {
	return c.InitTimeout
}

func (c ConfigApp) GetShutdownTimeout() *time.Duration {
	return c.ShutdownTimeout
}

type AppOption func(*App)

func WithLogger(logger appifaces.Logger) AppOption {
	return func(a *App) {
		a.logger = logger
	}
}

func WithProvidedSigs(sigs ...os.Signal) AppOption {
	return func(a *App) {
		a.shutdown.sigs = sigs
	}
}

func WithConfig(cfg appifaces.Configurator) AppOption {
	return func(a *App) {
		if cfg.IsSilientMode() {
			a.logger = logStub{}
		}
		if cfg.GetApplicationName() != nil {
			a.name = *cfg.GetApplicationName()
		}
		if cfg.GetInitTimeout() != nil {
			a.execution.initTimeout = cfg.GetInitTimeout()
		}
		if cfg.GetShutdownTimeout() != nil {
			a.shutdown.timeout = cfg.GetShutdownTimeout()
		}
	}
}

func (a *App) defaultSettingsCheckAndApply() {
	if a.name == "" {
		a.name = `unnamed`
	}

	if a.logger == nil {
		a.logger = defaultLog()
	}

	if a.shutdown.sigs == nil {
		a.shutdown.sigs = DefaultProvidedSigs
	}
	if a.shutdown.timeout == nil {
		a.shutdown.timeout = &DefaultShutdownTimeout
	}
}

// func WithName(name string) AppOption {
// 	return func(a *App) {
// 		a.name = name
// 	}
// }

// func WithInitTimeout(to time.Duration) AppOption {
// 	return func(a *App) {
// 		a.execution.initTimeout = &to
// 	}
// }

// func WithShutdownTimeout(to time.Duration) AppOption {
// 	return func(a *App) {
// 		a.shutdown.timeout = &to
// 	}
// }
