package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tusmasoma/go-microservice-k8s/microservice-k8s-demo/order/config"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type SQLExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

const (
	dbPrefix = "MYSQL_"
)

func NewMySQLDB(ctx context.Context) (*sql.DB, error) {
	conf, err := config.NewDBConfig(ctx, dbPrefix)
	if err != nil {
		log.Error("Failed to load database config", log.Ferror(err))
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Critical("Failed to connect to database", log.Fstring("dsn", dsn), log.Ferror(err))
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Critical("Failed to ping database", log.Ferror(err))
		return nil, err
	}

	log.Info("Successfully connected to database", log.Fstring("dsn", dsn))
	return db, nil
}
