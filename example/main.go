package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	ctapp "github.com/surkovvs/ct-app"
)

var (
	appName         string        = "example"
	initTimeout     time.Duration = time.Millisecond * 2500
	shutdownTimeout time.Duration = time.Millisecond * 2500
)

func main() {
	app := ctapp.New(
		ctapp.WithConfig(ctapp.ConfigApp{
			Silient:         false,
			Name:            &appName,
			InitTimeout:     &initTimeout,
			ShutdownTimeout: &shutdownTimeout,
		}),
		ctapp.WithLogger(slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				AddSource: false,
				Level:     slog.LevelDebug,
			}),
		)),
		ctapp.WithProvidedSigs(os.Interrupt),
	)

	app.AddModuleToGroup("example_group_1", "module_1", moduleInitRun{cfg: moduleCfg{
		Name: "mock_1",
		init: elemCfg{
			totalDur: time.Second * 2,
			wantFail: false,
		},
		run: elemCfg{
			totalDur: 0,
			wantFail: false,
		},
	}})

	app.AddModuleToGroup("example_group_2", "module_2", moduleInitRunSd{cfg: moduleCfg{
		Name: "mock_2",
		init: elemCfg{
			totalDur: time.Second * 2,
			wantFail: false,
		},
		run: elemCfg{
			totalDur: 0,
			wantFail: false,
		},
		shutdown: elemCfg{
			totalDur: time.Second * 2,
			wantFail: false,
		},
	}})

	ctx := context.Background()
	app.Start(ctx)
}
