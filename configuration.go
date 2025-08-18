package ctapp

import (
	"os"
	"time"

	"github.com/surkovvs/ct-app/appifaces"
)

type ConfigApp struct {
	Silient         bool
	Name            *string
	InitTimeout     *time.Duration
	ShutdownTimeout *time.Duration
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

type appOption func(*app)

func WithLogger(logger appifaces.Logger) appOption {
	return func(a *app) {
		a.logger = logger
	}
}

func WithProvidedSigs(sigs ...os.Signal) appOption {
	return func(a *app) {
		a.shutdown.sigs = sigs
	}
}

func WithConfig(cfg appifaces.Configurator) appOption {
	return func(a *app) {
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

func (a *app) defaultSettingsCheckAndApply() {
	if a.name == "" {
		a.name = `unnamed`
	}

	if a.logger == nil {
		a.logger = defaultLog()
	}

	if a.shutdown.sigs == nil {
		a.shutdown.sigs = defaultProvidedSigs
	}
	if a.shutdown.timeout == nil {
		a.shutdown.timeout = &defaultShutdownTimeout
	}
}

// func WithName(name string) appOption {
// 	return func(a *app) {
// 		a.name = name
// 	}
// }

// func WithInitTimeout(to time.Duration) appOption {
// 	return func(a *app) {
// 		a.execution.initTimeout = &to
// 	}
// }

// func WithShutdownTimeout(to time.Duration) appOption {
// 	return func(a *app) {
// 		a.shutdown.timeout = &to
// 	}
// }
