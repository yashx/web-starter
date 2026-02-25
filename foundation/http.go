package foundation

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type SubRouter interface {
	BuildHandler() (string, http.Handler)
}

func (app *App) StartHttpServer(subRouters ...SubRouter) error {
	router := chi.NewRouter()

	for _, sr := range subRouters {
		router.Mount(sr.BuildHandler())
	}

	port := app.Config.MustString("http.port")
	app.Logger.Info("Starting server", zap.String("port", port))

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		return err
	}

	return nil
}
