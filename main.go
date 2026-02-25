package main

import (
	"fmt"
	"web-starter/foundation"
	"web-starter/task"

	"go.uber.org/zap"
)

var (
	Version   = "unknown"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	app, err := foundation.InitApp()

	if err != nil {
		fmt.Println("Error initialising app")
		return
	}

	defer func(app *foundation.App) {
		err := app.Shutdown()
		if err != nil {
			app.Logger.Warn("Graceful shutdown failed", zap.Error(err))
		}
	}(app)

	defer func() {
		if r := recover(); r != nil {
			app.Logger.Sugar().Errorw("Application panicked", "panic", r)
		}
	}()

	app.Logger.Info("Starting Application",
		zap.String("version", Version),
		zap.String("build_time", BuildTime),
		zap.String("git_commit", GitCommit),
	)

	err = app.StartHttpServer(task.NewSubRouter(app))
	if err != nil {
		app.Logger.Panic("Failed to start http server", zap.Error(err))
	}
}
