package foundation

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type SubRouter interface {
	BuildHandler() (string, http.Handler)
}

func (app *App) StartHttpServer(ctx context.Context, subRouters ...SubRouter) error {
	router := chi.NewRouter()

	for _, sr := range subRouters {
		router.Mount(sr.BuildHandler())
	}

	port := app.Config.MustString("http.port")
	server := &http.Server{Addr: ":" + port, Handler: router}

	app.Logger.Info("starting server", zap.String("port", port))

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		app.Logger.Info("shutting down server")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
