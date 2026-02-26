package store

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go-api-scaffold/internal/model"
	"go-api-scaffold/pkg/config"
	"go-api-scaffold/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// Store is the data persistence layer
type Store struct {
	db *gorm.DB
}

// New creates a database connection
func New(cfg *config.DatabaseConfig) (*Store, error) {
	var dialector gorm.Dialector

	switch cfg.Type {
	case "sqlite":
		if dir := filepath.Dir(cfg.Path); dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create database directory: %w", err)
			}
		}
		dialector = sqlite.Open(cfg.Path)
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
		dialector = mysql.Open(dsn)
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)
		dialector = postgres.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// SQLite optimization
	if cfg.Type == "sqlite" {
		db.Exec("PRAGMA journal_mode=WAL")
		db.Exec("PRAGMA busy_timeout=5000")
		db.Exec("PRAGMA foreign_keys=ON")
	}

	// Connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	s := &Store{db: db}

	// Auto migrate
	if cfg.AutoMigrate {
		if err := s.AutoMigrate(); err != nil {
			return nil, fmt.Errorf("auto migrate: %w", err)
		}
		logger.Info("database migration completed")
	}

	logger.Infof("database connected: %s", cfg.Type)
	return s, nil
}

// DB returns the underlying GORM instance (for repositories)
func (s *Store) DB() *gorm.DB {
	return s.db
}

// Close closes the database connection
func (s *Store) Close() {
	if sqlDB, err := s.db.DB(); err == nil {
		_ = sqlDB.Close()
	}
}

// AutoMigrate runs auto migration for all models
func (s *Store) AutoMigrate() error {
	return s.db.AutoMigrate(
		&model.User{},
		&model.Example{},
		// GEN:MODEL_MIGRATE - Auto-appended by code generator, do not remove
	)
}
