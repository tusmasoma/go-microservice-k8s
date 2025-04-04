package mysql

import (
	"database/sql"
	"testing"

	"github.com/tusmasoma/go-microservice-k8s/go/pkg/log"
	"github.com/tusmasoma/go-microservice-k8s/go/pkg/mysql"

	_ "github.com/go-sql-driver/mysql" // This blank import is used for its init function
)

var (
	db        *sql.DB
	mysqlPort string
)

func TestMain(m *testing.M) {
	var closeMySQL func()
	var err error
	db, mysqlPort, closeMySQL, err = mysql.StartMySQL(db, "order")
	defer closeMySQL()
	if err != nil {
		log.Error("Failed to start MySQL: %v", err)
	}
	queries := []string{
		"DROP TABLE IF EXISTS OrderLines;",
		"DROP TABLE IF EXISTS Orders;",
	}
	if err = mysql.InitTestDatabase(db, queries); err != nil {
		return
	}
	m.Run()
}
