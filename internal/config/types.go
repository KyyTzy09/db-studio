package config

import "time"

// DriverType defines supported database types
type DriverType string

const (
	DriverPostgres DriverType = "postgres"
	DriverMySQL    DriverType = "mysql"
	DriverSQLite   DriverType = "sqlite"
)

// ConnectionConfig holds connection parameters for a database
type ConnectionConfig struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Driver      DriverType `json:"driver"`
	Host        string     `json:"host,omitempty"`
	Port        int        `json:"port,omitempty"`
	User        string     `json:"user,omitempty"`
	Password    string     `json:"password,omitempty"`
	Database    string     `json:"database,omitempty"`
	SSLMode     string     `json:"ssl_mode,omitempty"`
	FilePath    string     `json:"file_path,omitempty"` // For SQLite
	ProjectPath string     `json:"project_path"`        // Absolute path to project directory
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// GlobalConfig represents the root JSON stored in ~/.config/dbstudio/connections.json
type GlobalConfig struct {
	ActiveConnectionID string             `json:"active_connection_id,omitempty"`
	Connections        []ConnectionConfig `json:"connections"`
}
