package database

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/v2"
)

func GetConnection(config *koanf.Koanf) (*sqlx.DB, error) {
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
		return nil, err
	}

	db.SetMaxOpenConns(config.MustInt("database.pool.max-open"))
	db.SetMaxIdleConns(config.MustInt("database.pool.max-idle"))
	db.SetConnMaxLifetime(time.Duration(config.MustInt("database.pool.lifetime-minutes")) * time.Minute)
	db.SetConnMaxIdleTime(time.Duration(config.MustInt("database.pool.idle-minutes")) * time.Minute)

	return db, nil
}
