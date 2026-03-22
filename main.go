package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"web-starter/foundation"
	"web-starter/task"
	"web-starter/ws"

	"go.uber.org/zap"
)

var (
	Version   = "unknown"
	BuildTime = "unknown"
)

func main() {
	app, err := foundation.InitApp()

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error initialising app:", err)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			app.Logger.Sugar().Errorw("application panicked", "panic", r)
		}
	}()

	defer func(app *foundation.App) {
		err := app.Shutdown()
		if err != nil {
			app.Logger.Warn("graceful shutdown failed", zap.Error(err))
		}
	}(app)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app.Logger.Info("starting application",
		zap.String("version", Version),
		zap.String("build_time", BuildTime),
	)

	err = app.StartHttpServer(ctx, task.NewSubRouter(app), ws.NewSubRouter(app))
	if err != nil {
		app.Logger.Error("http server stopped", zap.Error(err))
	}
}
