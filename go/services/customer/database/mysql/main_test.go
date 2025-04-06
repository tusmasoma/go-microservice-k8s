package mysql

import (
	"context"
	"database/sql"
	"testing"

	"github.com/tusmasoma/go-microservice-k8s/go/pkg/log"
	pt "github.com/tusmasoma/go-microservice-k8s/go/pkg/testing"

	_ "github.com/go-sql-driver/mysql" // This blank import is used for its init function
)

var db *sql.DB

func TestMain(m *testing.M) {
	var closeMySQL func()
	var err error
	ctx := context.Background()
	params := pt.MySQLParams{
		DB:           db,
		DBName:       "go-microservice-k8s-test",
		ServiceName:  "customer",
		UserName:     "root",
		UserPassword: "root",
	}
	db, closeMySQL, err = pt.StartMySQL(ctx, &params)
	defer closeMySQL()
	if err != nil {
		log.Error("Failed to start MySQL: %v", err)
	}
	queries := []string{
		"DROP TABLE IF EXISTS Customers;",
	}
	if err = pt.InitTestDatabase(ctx, &params, queries); err != nil {
		return
	}
	m.Run()
}
