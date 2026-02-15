package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	GRPC     GRPCConfig     `mapstructure:"grpc"`
	Log      LogConfig      `mapstructure:"log"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Mode    string `mapstructure:"mode"` // debug, release, test
}

type ServerConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout"`  // seconds
	WriteTimeout int    `mapstructure:"write_timeout"` // seconds
}

type GRPCConfig struct {
	Enabled bool `mapstructure:"enabled"`
	Port    int  `mapstructure:"port"`
}

type DatabaseConfig struct {
	Type            string `mapstructure:"type"` // sqlite, mysql, postgres
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	Database        string `mapstructure:"database"`
	Path            string `mapstructure:"path"` // SQLite file path
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"` // minutes
	AutoMigrate     bool   `mapstructure:"auto_migrate"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`       // debug, info, warn, error
	Format     string `mapstructure:"format"`      // console, json
	Output     string `mapstructure:"output"`      // console, file, both
	FilePath   string `mapstructure:"file_path"`   // log file path
	MaxSize    int    `mapstructure:"max_size"`    // MB per file
	MaxBackups int    `mapstructure:"max_backups"` // rotation count
	MaxAge     int    `mapstructure:"max_age"`     // days retention
	Compress   bool   `mapstructure:"compress"`    // gzip old logs
}

type JWTConfig struct {
	Secret       string `mapstructure:"secret"`
	Expire       int    `mapstructure:"expire"`        // hours
	RefreshHours int    `mapstructure:"refresh_hours"` // refresh window in hours
}

// Load reads configuration from file
func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	// Support environment variable override with APP_ prefix
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	cfg := DefaultConfig()
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return cfg, nil
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:    "my-service",
			Version: "1.0.0",
			Mode:    "debug",
		},
		Server: ServerConfig{
			Host:         "0.0.0.0",
			Port:         8080,
			ReadTimeout:  10,
			WriteTimeout: 10,
		},
		GRPC: GRPCConfig{
			Enabled: false,
			Port:    9090,
		},
		Database: DatabaseConfig{
			Type:            "sqlite",
			Path:            "./data/app.db",
			MaxOpenConns:    100,
			MaxIdleConns:    10,
			ConnMaxLifetime: 60,
			AutoMigrate:     true,
		},
		Log: LogConfig{
			Level:      "info",
			Format:     "console",
			Output:     "both",
			FilePath:   "logs/app.log",
			MaxSize:    50,
			MaxBackups: 3,
			MaxAge:     7,
			Compress:   true,
		},
		JWT: JWTConfig{
			Secret:       "change-me-in-production",
			Expire:       24,
			RefreshHours: 168, // 7 days
		},
	}
}

// Validate checks configuration validity
func (c *Config) Validate() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	switch c.Database.Type {
	case "sqlite":
		if c.Database.Path == "" {
			return fmt.Errorf("sqlite requires database.path")
		}
	case "mysql", "postgres":
		if c.Database.Host == "" || c.Database.Database == "" {
			return fmt.Errorf("%s requires database.host and database.database", c.Database.Type)
		}
	default:
		return fmt.Errorf("unsupported database type: %s", c.Database.Type)
	}

	if c.JWT.Secret == "" {
		return fmt.Errorf("jwt.secret is required")
	}

	return nil
}
