package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"
)

type totaller interface {
	Healthcheck(ctx context.Context) error
	Init(ctx context.Context) error
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

var customError = errors.New("custom error")

const (
	healthcheckStarted = `Healthcheck started`
	healthcheckDone    = `Healthcheck done`

	initStarted = `Init started`
	initDone    = `Init done`

	runStarted = `Run started`
	runRunning = `Run in process`
	runDone    = `Run done`

	shutdownStarted = `Shutdown started`
	shutdownDone    = `Shutdown done`
)

type elemCfg struct {
	totalDur time.Duration
	wantFail bool
}

type moduleCfg struct {
	Name        string
	healthcheck elemCfg
	init        elemCfg
	run         elemCfg
	shutdown    elemCfg
}

func Executing(ctx context.Context, module any, name string, cfg elemCfg, stage string) error {
	exCtx := context.Background()
	if cfg.totalDur != 0 {
		exCtx, _ = context.WithTimeout(ctx, cfg.totalDur)
	}
	ticker := time.NewTicker(time.Second * 1)
EndlessCycle:
	for {
		select {
		case <-ctx.Done():
			log.Printf("%s [%s]: %s\n", reflect.ValueOf(module).Type().Name(), name, "stage "+stage+" ended by external ctx")
			return customError
			// break EndlessCycle
		case <-exCtx.Done():
			break EndlessCycle
		case <-ticker.C:
			log.Printf("%s [%s]: %s\n", reflect.ValueOf(module).Type().Name(), name, "executing "+stage)
		}
	}
	if cfg.wantFail {
		return fmt.Errorf("%s [%s]:%w", reflect.ValueOf(module).Type().Name(), name, customError)
	}
	return nil
}

type moduleInitRun struct {
	cfg moduleCfg
}

func (m moduleInitRun) Init(ctx context.Context) error {
	log.Printf("%s: %s\n", reflect.ValueOf(m).Type().Name(), initStarted)
	defer log.Printf("%s: %s\n", reflect.ValueOf(m).Type().Name(), initDone)
	return Executing(ctx, m, m.cfg.Name, m.cfg.init, "init")
}

func (m moduleInitRun) Run(ctx context.Context) error {
	log.Printf("%s: %s\n", reflect.ValueOf(m).Type().Name(), runStarted)
	defer log.Printf("%s: %s\n", reflect.ValueOf(m).Type().Name(), runDone)
	return Executing(ctx, m, m.cfg.Name, m.cfg.run, "run")
}

type moduleInitRunSd struct {
	cfg moduleCfg
}

func (m moduleInitRunSd) Init(ctx context.Context) error {
	log.Printf("%s: %s\n", reflect.ValueOf(m).Type().Name(), initStarted)
	defer log.Printf("%s: %s\n", reflect.ValueOf(m).Type().Name(), initDone)
	return Executing(ctx, m, m.cfg.Name, m.cfg.init, "init")
}

func (m moduleInitRunSd) Run(ctx context.Context) error {
	log.Printf("%s: %s\n", reflect.ValueOf(m).Type().Name(), runStarted)
	defer log.Printf("%s: %s\n", reflect.ValueOf(m).Type().Name(), runDone)
	return Executing(ctx, m, m.cfg.Name, m.cfg.run, "run")
}

func (m moduleInitRunSd) Shutdown(ctx context.Context) error {
	log.Printf("%s: %s\n", reflect.ValueOf(m).Type().Name(), shutdownStarted)
	defer log.Printf("%s: %s\n", reflect.ValueOf(m).Type().Name(), shutdownDone)
	return Executing(ctx, m, m.cfg.Name, m.cfg.shutdown, "shutdown")
}

type moduleHchRunSd struct {
	cfg moduleCfg
}

func (m moduleHchRunSd) Healthcheck(ctx context.Context) error {
	log.Printf("%s: %s\n", reflect.ValueOf(m).Type().Name(), initStarted)
	defer log.Printf("%s: %s\n", reflect.ValueOf(m).Type().Name(), initDone)
	return Executing(ctx, m, m.cfg.Name, m.cfg.init, "init")
}

func (m moduleHchRunSd) Run(ctx context.Context) error {
	log.Printf("%s: %s\n", reflect.ValueOf(m).Type().Name(), runStarted)
	defer log.Printf("%s: %s\n", reflect.ValueOf(m).Type().Name(), runDone)
	return Executing(ctx, m, m.cfg.Name, m.cfg.run, "run")
}

func (m moduleHchRunSd) Shutdown(ctx context.Context) error {
	log.Printf("%s: %s\n", reflect.ValueOf(m).Type().Name(), shutdownStarted)
	defer log.Printf("%s: %s\n", reflect.ValueOf(m).Type().Name(), shutdownDone)
	return Executing(ctx, m, m.cfg.Name, m.cfg.shutdown, "shutdown")
}

type moduleSd struct {
	cfg moduleCfg
}

func (m moduleSd) Shutdown(ctx context.Context) error {
	log.Printf("%s: %s\n", reflect.ValueOf(m).Type().Name(), shutdownStarted)
	defer log.Printf("%s: %s\n", reflect.ValueOf(m).Type().Name(), shutdownDone)
	return Executing(ctx, m, m.cfg.Name, m.cfg.shutdown, "shutdown")
}
