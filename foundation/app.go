package foundation

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"web-starter/foundation/internal/config"
	"web-starter/foundation/internal/database"

	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

type App struct {
	Logger *zap.Logger
	DB     *sqlx.DB
	Config *koanf.Koanf
}

func InitApp() (*App, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error initializing logger:", err)
		return nil, err
	}

	appConfig := config.GetConfig(logger)

	db, err := database.GetConnection(appConfig)
	if err != nil {
		return nil, err
	}

	if err = database.RunMigrations(logger, db.DB); err != nil {
		return nil, err
	}

	app := App{
		Logger: logger,
		DB:     db,
		Config: appConfig,
	}

	return &app, nil
}

func (app *App) Shutdown() error {
	err := app.DB.Close()
	if err != nil {
		return err
	}

	err = app.Logger.Sync()
	if err != nil && !errors.Is(err, syscall.EINVAL) {
		return err
	}

	return nil
}
