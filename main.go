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
)

func main() {
	app, err := foundation.InitApp()

	if err != nil {
		fmt.Println("error initialising app")
		return
	}

	defer func(app *foundation.App) {
		err := app.Shutdown()
		if err != nil {
			app.Logger.Warn("graceful shutdown failed", zap.Error(err))
		}
	}(app)

	defer func() {
		if r := recover(); r != nil {
			app.Logger.Sugar().Errorw("application panicked", "panic", r)
		}
	}()

	app.Logger.Info("starting application",
		zap.String("version", Version),
		zap.String("build_time", BuildTime),
	)

	err = app.StartHttpServer(task.NewSubRouter(app))
	if err != nil {
		app.Logger.Panic("failed to start http server", zap.Error(err))
	}
}
