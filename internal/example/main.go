package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	ctapp "github.com/surkovvs/ct-app"
	"github.com/surkovvs/ct-app/internal/example/modules"
)

var (
	appName         string        = "example"
	initTimeout     time.Duration = time.Millisecond * 2000
	shutdownTimeout time.Duration = time.Millisecond * 2000
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
				Level: slog.LevelDebug,
			}),
		)),
		ctapp.WithProvidedSigs(os.Interrupt),
	)

	app.AddModuleToGroup("example_group_1", "module_1:1", modules.NewModuleHcInitRun(modules.ModuleHcInitRunCfg{
		Name: "mock_1:1",
		Init: modules.ElemCfg{
			TotalDur: time.Millisecond * 2100,
			WantFail: false,
		},
	}))

	app.AddModuleToGroup("example_group_1", "module_1:2", modules.NewModuleInitRunSd(modules.ModuleInitRunSdCfg{
		Name: "mock_1:2",
		Init: modules.ElemCfg{
			TotalDur: time.Millisecond * 1600,
			WantFail: false,
		},
		Shutdown: modules.ElemCfg{
			TotalDur: time.Millisecond * 1600,
			WantFail: false,
		},
	}))

	app.AddModuleToGroup("example_group_2", "module_2", modules.NewModuleRunSd(modules.ModuleRunSdCfg{
		Name: "mock_2",
		Run: modules.ElemCfg{
			TotalDur: time.Millisecond * 4000,
			WantFail: false,
		},
		Shutdown: modules.ElemCfg{
			TotalDur: time.Millisecond * 600,
			WantFail: false,
		},
	}))

	app.AddBackgroundModule("module_bg_1", modules.NewModuleInitRunSd(modules.ModuleInitRunSdCfg{
		Name: "mock_bg_2",
		Init: modules.ElemCfg{
			TotalDur: time.Millisecond * 1000,
			WantFail: false,
		},
		Run: modules.ElemCfg{
			TotalDur: 3000,
			WantFail: false,
		},
		Shutdown: modules.ElemCfg{
			TotalDur: time.Millisecond * 1000,
			WantFail: false,
		},
	}))

	app.AddBackgroundModule("module_bg_2", modules.NewModuleInitSd(modules.ModuleInitSdCfg{
		Name: "mock_bg_2",
		Init: modules.ElemCfg{
			TotalDur: time.Millisecond * 1000,
			WantFail: false,
		},
		Shutdown: modules.ElemCfg{
			TotalDur: time.Millisecond * 1000,
			WantFail: false,
		},
	}))

	ctx := context.Background()
	app.Start(ctx)
}
