package foundation

import (
	"errors"
	"fmt"
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
		fmt.Println("Error initializing logger")
		return nil, err
	}

	appConfig := config.GetConfig(logger)

	db := database.GetConnection(logger, appConfig)
	err = database.RunMigrations(logger, db.DB)
	if err != nil {
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
	err := app.Logger.Sync()
	if err != nil && !errors.Is(err, syscall.EINVAL) {
		return err
	}

	err = app.DB.Close()
	if err != nil {
		return err
	}

	return nil
}
