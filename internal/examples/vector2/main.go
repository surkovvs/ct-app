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
)

const (
	init3Duration     = time.Millisecond * 250
	run1Duration      = time.Millisecond * 500
	shutdown1Duration = time.Millisecond * 500

	releaseShutdown2Since = time.Millisecond * 1250
	releaseGroup3Since    = time.Millisecond * 1000

	cancelExecCtxSince   = time.Millisecond * 250
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
	ctx := context.Background()

	stubWaiting := make(chan struct{})

	c := vector.NewConstructor(logReceiver())
	v := c.Sequentially(ctx, stub("init_1", 0), stub("init_2", 0), stub("init_3", 0),
		c.WithRelease(stubWaiting, stub("stub", 0)),
	)

	fmt.Printf("vector builded since: %v\n", time.Since(start))

	go func() {
		time.Sleep(cancelExecCtxSince)
		cancelExecCtx()
	}()

	execution := time.Now()

	go func() {
		<-stubWaiting
		fmt.Printf("stubWaiting released since vector build: %v\n", time.Since(execution))
	}()

	v.Exec(execCtx)
	<-execCtx.Done()
	defer time.Sleep(time.Millisecond * 2000)

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
