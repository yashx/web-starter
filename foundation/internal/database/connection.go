package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

func GetConnection(logger *zap.Logger, config *koanf.Koanf) *sqlx.DB {
	username := config.MustString("database.username")
	password := config.MustString("database.password")
	host := config.MustString("database.host")
	port := config.MustString("database.port")
	name := config.MustString("database.name")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		username,
		password,
		host,
		port,
		name,
	)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		logger.Panic("failed to connect to database: %v", zap.Error(err))
	}
	return db
}
