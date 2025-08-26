package ctapp

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/surkovvs/ct-app/appifaces"
	"github.com/surkovvs/ct-app/compstor"
)

var (
	defaultProvidedSigs    = []os.Signal{syscall.SIGTERM, syscall.SIGINT}
	defaultShutdownTimeout = time.Second * 3
	BackgroundGroup        = `background`
)

type (
	execution struct {
		wg            *sync.WaitGroup
		done          chan struct{}
		errFlow       chan error
		initRunCancel context.CancelFunc
		initTimeout   *time.Duration
	}
	shutdown struct {
		ctx          context.Context
		ctxCancel    context.CancelFunc
		wg           *sync.WaitGroup
		shutdownDone chan struct{}
		sigs         []os.Signal
		timeout      *time.Duration
		exitCode     int
	}
	app struct {
		execution execution
		shutdown  shutdown
		storage   compstor.CompsStorage
		name      string
		logger    appifaces.Logger
	}
)

func New(opts ...appOption) *app {
	sdCtx, sdCancel := context.WithCancel(context.Background())
	a := &app{
		execution: execution{
			wg:            &sync.WaitGroup{},
			done:          make(chan struct{}),
			errFlow:       make(chan error),
			initRunCancel: nil,
			initTimeout:   nil,
		},
		shutdown: shutdown{
			ctx:          sdCtx,
			ctxCancel:    sdCancel,
			wg:           &sync.WaitGroup{},
			shutdownDone: make(chan struct{}),
			sigs:         nil,
			timeout:      nil,
			exitCode:     0,
		},
		storage: compstor.NewCompsStorage(),
		name:    "",
		logger:  nil,
	}

	for _, opt := range opts {
		opt(a)
	}
	a.defaultSettingsCheckAndApply()

	go a.accompaniment()

	if err := a.storage.AddGroup(BackgroundGroup); err != nil &&
		!errors.Is(err, compstor.ErrGroupAlreadyRegistered) {
		log.Fatal(err)
	}

	return a
}

func (a *app) accompaniment() {
	syscallC := make(chan os.Signal, 1)
	signal.Notify(syscallC, a.shutdown.sigs...)
	for {
		select {
		case err := <-a.execution.errFlow:
			a.logger.Error(`module error`,
				"application", a.name,
				`error`, err)
		case <-a.execution.done:
			a.execution.done = nil
			a.logger.Debug(`execution finished graceful shutdown started`,
				"application", a.name)
			signal.Stop(syscallC)
			a.execution.initRunCancel()
			go a.gracefulShutdown()
		case sig := <-syscallC:
			a.logger.Info(`graceful shutdown started by syscall`,
				"application", a.name,
				`syscall`, sig.String())

			signal.Stop(syscallC)
			a.execution.initRunCancel()
			go a.gracefulShutdown()
		}
	}
}
