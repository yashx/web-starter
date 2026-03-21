package database

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func RunMigrations(logger *zap.Logger, DB *sql.DB) error {
	gl := zapGooseLogger{
		logger.Sugar(),
	}
	goose.SetLogger(&gl)
	goose.SetBaseFS(embedMigrations)
	err := goose.SetDialect("mysql")
	if err != nil {
		return err
	}

	return goose.Up(DB, "migrations")
}

type zapGooseLogger struct {
	log *zap.SugaredLogger
}

func (l *zapGooseLogger) Printf(format string, v ...interface{}) {
	l.log.Infof(format, v...)
}

func (l *zapGooseLogger) Fatalf(format string, v ...interface{}) {
	l.log.Fatalf(format, v...)
}
