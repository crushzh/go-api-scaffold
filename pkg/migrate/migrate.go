// Package migrate 提供基于 golang-migrate 的数据库迁移支持。
// 生产模式使用 SQL 迁移文件，开发模式使用 GORM AutoMigrate。
package migrate

import (
	"fmt"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	// database drivers
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
)

// Config 迁移配置
type Config struct {
	Driver     string // sqlite, mysql, postgres
	Host       string
	Port       int
	User       string
	Password   string
	DBName     string
	SQLitePath string
}

// Run 执行数据库迁移
func Run(migrationFS fs.FS, dir string, cfg *Config) error {
	source, err := iofs.New(migrationFS, dir)
	if err != nil {
		return fmt.Errorf("create migration source: %w", err)
	}

	dbURL, err := buildURL(cfg)
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, dbURL)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration up: %w", err)
	}
	return nil
}

func buildURL(cfg *Config) (string, error) {
	switch cfg.Driver {
	case "sqlite":
		path := cfg.SQLitePath
		if path == "" {
			path = "data/app.db"
		}
		return "sqlite://" + path, nil
	case "mysql":
		return fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s?multiStatements=true",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName), nil
	case "postgres":
		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName), nil
	default:
		return "", fmt.Errorf("unsupported driver: %s", cfg.Driver)
	}
}
