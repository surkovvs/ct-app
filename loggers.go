package ctapp

import (
	"log/slog"
	"os"

	"github.com/surkovvs/ct-app/appifaces"
)

type logStub struct{}

func (logStub) Debug(string, ...any) {}

func (logStub) Info(string, ...any) {}

func (logStub) Warn(string, ...any) {}

func (logStub) Error(string, ...any) {}

type logWrap struct {
	logger appifaces.Logger
}

func defaultLog() logWrap {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	return logWrap{
		logger: logger,
	}
}

func (l logWrap) Debug(msg string, args ...any) {
	l.logger.Debug(`[GoCAT] `+msg, args...)
}

func (l logWrap) Info(msg string, args ...any) {
	l.logger.Info(`[GoCAT] `+msg, args...)
}

func (l logWrap) Warn(msg string, args ...any) {
	l.logger.Warn(`[GoCAT] `+msg, args...)
}

func (l logWrap) Error(msg string, args ...any) {
	l.logger.Error(`[GoCAT] `+msg, args...)
}
