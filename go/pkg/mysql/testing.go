package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/tusmasoma/go-microservice-k8s/go/pkg/log"
)

func StartMySQL(db *sql.DB, serviceName string) (*sql.DB, string, func(), error) {
	migrationsDir, err := findMigrationsDir()
	if err != nil {
		log.Error("Failed to find migrations directory: %v", err)
		return nil, "", nil, err
	}
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Error("Could not connect to Docker: %s", err)
		return nil, "", nil, err
	}
	err = pool.Client.Ping()
	if err != nil {
		log.Error("Could not ping Docker: %s", err)
		return nil, "", nil, err
	}
	runOptions := &dockertest.RunOptions{
		Repository: "mysql",
		Tag:        "8.0",
		Env: []string{
			"MYSQL_ROOT_USER=root",
			"MYSQL_ROOT_PASSWORD=microservice-k8s-demo",
			"MYSQL_DATABASE=microservice-k8s-demo-test-db",
		},
		Cmd: []string{
			"--character-set-server=utf8mb4",
			"--collation-server=utf8mb4_unicode_ci",
		},
	}
	resource, err := pool.RunWithOptions(runOptions,
		func(hc *docker.HostConfig) {
			hc.AutoRemove = true
			hc.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
			hc.Mounts = []docker.HostMount{
				{
					Type:   "bind",
					Source: filepath.Join(migrationsDir, serviceName, serviceName+".ddl"),
					Target: fmt.Sprintf("/docker-entrypoint-initdb.d/%s.ddl", serviceName),
				},
			}
		},
	)
	if err != nil {
		log.Error("Could not start resource: %s", err)
		return nil, "", nil, err
	}
	port := resource.GetPort("3306/tcp")
	err = pool.Retry(func() error {
		dsn := fmt.Sprintf("root:microservice-k8s-demo@(localhost:%s)/microservice-k8s-demo-test-db?charset=utf8mb4&parseTime=true", port)
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			return err
		}
		return db.Ping()
	})
	if err != nil {
		log.Error("Could not connect to MySQL: %s", err)
		return nil, "", nil, err
	}
	log.Info("start MySQL container🐳")
	return db, port, func() { closeMySQL(db, pool, resource) }, nil
}

func closeMySQL(db *sql.DB, pool *dockertest.Pool, resource *dockertest.Resource) {
	if err := db.Close(); err != nil {
		log.Error("Failed to close MySQL connection: %v", err)
	}
	if err := pool.Purge(resource); err != nil {
		log.Error("Failed to purge MySQL container: %v", err)
	}
	log.Info("close MySQL container🐳")
}

func findMigrationsDir() (string, error) {
	const migrationsDirName = "migrations"
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}
	for {
		migrationsPath := filepath.Join(dir, migrationsDirName)
		info, err := os.Stat(migrationsPath)
		if err == nil && info.IsDir() {
			return migrationsPath, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("migrations directory not found in any parent")
		}
		dir = parent
	}
}

func InitTestDatabase(db *sql.DB, queries []string) error {
	initDBQuery := `
	CREATE DATABASE IF NOT EXISTS microservice-k8s-demo-test-db DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
	USE microservice-k8s-demo-test-db;
	`
	initQuery := initDBQuery + strings.Join(queries, "\n")
	if _, err := db.Exec(initQuery); err != nil {
		return fmt.Errorf("failed to initialize test DB: %w", err)
	}
	return nil
}

func ValidateErr(t *testing.T, err error, wantErr error) {
	if (err != nil) != (wantErr != nil) {
		t.Errorf("error = %v, wantErr %v", err, wantErr)
	} else if err != nil && wantErr != nil && err.Error() != wantErr.Error() {
		t.Errorf("error = %v, wantErr %v", err, wantErr)
	}
}
