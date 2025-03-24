package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/go-microservice-k8s/pkg/config"
)

type DB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
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
		log.Critical("Failed to connect to database", log.Ferror(err))
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Critical("Failed to ping database", log.Ferror(err))
		return nil, err
	}

	log.Info("Successfully connected to database")
	return db, nil
}
