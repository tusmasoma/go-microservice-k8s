package testing

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/tusmasoma/go-microservice-k8s/go/pkg/log"
)

type MySQLParams struct {
	DB           *sql.DB
	DBName       string
	UserName     string
	UserPassword string
}

func StartMySQL(ctx context.Context, params *MySQLParams) (*sql.DB, func(), error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Error("Could not connect to Docker: %s", err)
		return nil, nil, err
	}
	err = pool.Client.PingWithContext(ctx)
	if err != nil {
		log.Error("Could not ping Docker: %s", err)
		return nil, nil, err
	}
	runOptions := &dockertest.RunOptions{
		Repository: "mysql",
		Tag:        "8.0",
		Env: []string{
			fmt.Sprintf("MYSQL_ROOT_USER=%s", params.UserName),
			fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", params.UserPassword),
			fmt.Sprintf("MYSQL_DATABASE=%s", params.DBName),
		},
		Cmd: []string{
			"--character-set-server=utf8mb4",
			"--collation-server=utf8mb4_unicode_ci",
		},
	}
	// migrationsDir, err := findMigrationsDir()
	// if err != nil {
	// 	log.Error("Failed to find migrations directory: %v", err)
	// 	return nil, nil, err
	// }
	resource, err := pool.RunWithOptions(runOptions,
		func(hc *docker.HostConfig) {
			hc.AutoRemove = true
			hc.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
			hc.Mounts = []docker.HostMount{
				// {
				// 	Type:   "bind",
				// 	Source: filepath.Join(migrationsDir, params.ServiceName, params.ServiceName+".ddl"),
				// 	Target: fmt.Sprintf("/docker-entrypoint-initdb.d/%s.ddl", params.ServiceName),
				// },
			}
		},
	)
	if err != nil {
		log.Error("Could not start resource: %s", err)
		return nil, nil, err
	}
	port := resource.GetPort("3306/tcp")
	err = pool.Retry(func() error {
		dsn := fmt.Sprintf("%s:%s@(localhost:%s)/%s?charset=utf8mb4&parseTime=true",
			params.UserName, params.UserPassword, port, params.DBName)
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return err
		}
		if db == nil {
			log.Error("Database connection is nil")
			return fmt.Errorf("database connection is nil")
		}
		if err := db.PingContext(ctx); err != nil {
			return err
		}
		params.DB = db
		return nil
	})
	if err != nil {
		log.Error("Could not connect to MySQL: %s", err)
		return nil, nil, err
	}
	log.Info("start MySQL container🐳")
	return params.DB, func() { closeMySQL(params.DB, pool, resource) }, nil
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

func InitTestDatabase(ctx context.Context, params *MySQLParams, queries []string) error {
	initDBQuery := fmt.Sprintf(`
	CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
	USE %s;
	`, params.DBName, params.DBName)
	initQuery := initDBQuery + strings.Join(queries, "\n")
	if _, err := params.DB.ExecContext(ctx, initQuery); err != nil {
		return fmt.Errorf("failed to initialize test DB: %w", err)
	}
	return nil
}
