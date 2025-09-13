package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/surkovvs/ct-app/vector"
	"github.com/surkovvs/ct-app/wgchan"
)

const (
	init3Duration     = time.Millisecond * 250
	run1Duration      = time.Millisecond * 500
	shutdown1Duration = time.Millisecond * 500

	releaseShutdown2Since = time.Millisecond * 1250
	releaseGroup3Since    = time.Millisecond * 1000

	cancelExecCtxSince   = time.Millisecond * 1500
	cancelGroup1CtxSince = time.Millisecond * 750
)

func main() {
	start := time.Now()

	// log.SetFlags(log.Lmicroseconds)
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "since_start" {
				a.Value = slog.DurationValue(a.Value.Time().Sub(start))
			}
			return a
		},
	}),
	),
	)

	execCtx, cancelExecCtx := context.WithCancel(context.Background())
	group1Ctx, cancelGroup1Ctx := context.WithCancel(context.Background())
	ctx := context.Background()

	stubWaiting := make(chan struct{})
	initWaiting := make(chan struct{})
	blockShutdown2 := make(chan struct{})
	blockGroup3 := make(chan struct{})

	wgc := wgchan.NewWgChan()

	c := vector.NewConstructor(logReceiver())
	v := c.Sequentially(ctx,
		c.WithRelease(initWaiting,
			c.Concurrently(ctx,
				stub("init_1", 0), stub("init_2", 0), stub("init_3", init3Duration),
			),
		),
		c.Concurrently(ctx,
			c.Sequentially(group1Ctx,
				stub("run_1", run1Duration),
				c.WithRelease(wgc.AddChan(),
					stub("shutdown_1", shutdown1Duration),
				),
			),
			c.Sequentially(ctx,
				stub("run_2", 0),
				c.WithRelease(wgc.AddChan(),
					c.WithWait(blockShutdown2,
						stub("shutdown_2", 0),
					)),
			),
			c.WithWait(blockGroup3,
				c.Sequentially(ctx,
					stub("run_3", 0),
					c.WithRelease(wgc.AddChan(),
						stub("shutdown_3", 0)),
				),
			),
		),
		c.WithRelease(stubWaiting, stub("stub", 0)),
	)

	fmt.Printf("vector builded since: %v\n", time.Since(start))

	go func() {
		time.Sleep(releaseShutdown2Since)
		close(blockShutdown2)
	}()

	go func() {
		time.Sleep(releaseGroup3Since)
		close(blockGroup3)
	}()

	go func() {
		time.Sleep(cancelExecCtxSince)
		cancelExecCtx()
	}()

	go func() {
		time.Sleep(cancelGroup1CtxSince)
		cancelGroup1Ctx()
	}()

	execution := time.Now()

	go func() {
		<-stubWaiting
		fmt.Printf("stubWaiting released since vector build: %v\n", time.Since(execution))
	}()

	go func() {
		<-initWaiting
		fmt.Printf("initWaiting released since vector build: %v\n", time.Since(execution))
	}()

	v.Exec(execCtx)
	select {
	case <-wgc.WaitChan():
	case <-execCtx.Done():
		defer time.Sleep(time.Millisecond * 2000)
	}

	fmt.Printf("all shutdowns ended since: %v\n", time.Since(execution))
	fmt.Printf("total run time since start: %v\n", time.Since(start))
}

func logReceiver() chan<- error {
	errChan := make(chan error)
	go func() {
		for err := range errChan {
			log.Println("err", err)
		}
	}()

	return errChan
}

func stub(tag string, dur time.Duration) func(context.Context) error {
	return func(ctx context.Context) error {
		if dur != 0 {
			time.Sleep(dur)
		}
		if ctx.Err() != nil {
			return errors.New(tag + " ctx canceled")
		}
		slog.Info("stub report", "since_start", time.Now(), "tag", tag, "status", "finished")
		return nil
	}
}
