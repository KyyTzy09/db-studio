package db

import (
	"fmt"

	"db-studio-go/internal/config"
)

// NewDriver constructs a Database driver instance based on ConnectionConfig
func NewDriver(cfg config.ConnectionConfig) (Database, error) {
	switch cfg.Driver {
	case config.DriverPostgres:
		return NewPostgresDriver(cfg), nil
	case config.DriverMySQL:
		return NewMySQLDriver(cfg), nil
	case config.DriverSQLite:
		return NewSQLiteDriver(cfg), nil
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}
}
