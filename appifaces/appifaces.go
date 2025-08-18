package appifaces

import (
	"context"
	"time"
)

type (
	Healthchecker interface {
		Healthcheck(ctx context.Context) error
	}
	Initializer interface {
		Init(ctx context.Context) error
	}
	Runner interface {
		Run(ctx context.Context) error
	}
	Shutdowner interface {
		Shutdown(ctx context.Context) error
	}
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type Configurator interface {
	IsSilientMode() bool
	GetApplicationName() *string
	GetInitTimeout() *time.Duration
	GetShutdownTimeout() *time.Duration
}
