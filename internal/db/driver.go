package db

import (
	"context"

	"db-studio-go/internal/config"
)

// ColumnInfo holds column metadata
type ColumnInfo struct {
	Name            string `json:"name"`
	DataType        string `json:"data_type"`
	IsNullable      bool   `json:"is_nullable"`
	IsPrimaryKey    bool   `json:"is_primary_key"`
	IsForeignKey    bool   `json:"is_foreign_key,omitempty"`
	IsAutoIncrement bool   `json:"is_auto_increment,omitempty"`
	DefaultValue    string `json:"default_value,omitempty"`
}

// TableInfo holds high-level table description
type TableInfo struct {
	Name     string `json:"name"`
	Type     string `json:"type"` // BASE TABLE, VIEW
	RowCount int64  `json:"row_count,omitempty"`
}

// TableSchema holds detailed schema information
type TableSchema struct {
	TableName string       `json:"table_name"`
	Columns   []ColumnInfo `json:"columns"`
}

// QueryResult holds execution response for data fetch or raw queries
type QueryResult struct {
	Columns      []string                 `json:"columns"`
	Rows         []map[string]interface{} `json:"rows"`
	AffectedRows int64                    `json:"affected_rows"`
	ExecutionMs  int64                    `json:"execution_ms"`
}

// ColumnSpec defines a column for table creation and schema alteration
type ColumnSpec struct {
	Name             string `json:"name"`
	DataType         string `json:"data_type"`
	IsPrimaryKey     bool   `json:"is_primary_key"`
	IsNullable       bool   `json:"is_nullable"`
	DefaultValue     string `json:"default_value,omitempty"`
	AutoIncrement    bool   `json:"auto_increment,omitempty"`
	ForeignKeyTable  string `json:"fk_table,omitempty"`
	ForeignKeyColumn string `json:"fk_column,omitempty"`
}

// CreateTableRequest holds the payload to create a new table
type CreateTableRequest struct {
	TableName string       `json:"table_name"`
	Columns   []ColumnSpec `json:"columns"`
}

// ForeignKeyRelation holds a foreign key link between two tables
type ForeignKeyRelation struct {
	ID           string `json:"id"`
	SourceTable  string `json:"source_table"`
	SourceColumn string `json:"source_column"`
	TargetTable  string `json:"target_table"`
	TargetColumn string `json:"target_column"`
}

// SchemaGraph holds complete database nodes and FK relationships for ER Diagram
type SchemaGraph struct {
	Nodes []TableSchema        `json:"nodes"`
	Edges []ForeignKeyRelation `json:"edges"`
}

// Database defines standard operations for supported database drivers
type Database interface {
	Connect(ctx context.Context) error
	Disconnect() error
	Ping(ctx context.Context) error
	GetTables(ctx context.Context) ([]TableInfo, error)
	GetSchema(ctx context.Context, tableName string) (*TableSchema, error)
	GetSchemaGraph(ctx context.Context) (*SchemaGraph, error)
	ExecuteQuery(ctx context.Context, query string, force bool) (*QueryResult, error)
	InsertRow(ctx context.Context, table string, data map[string]interface{}) error
	UpdateRow(ctx context.Context, table string, pk map[string]interface{}, data map[string]interface{}) error
	DeleteRow(ctx context.Context, table string, pk map[string]interface{}) error
	BatchInsertOrUpdate(ctx context.Context, table string, rows []map[string]interface{}, mode string) (int64, error)
	CreateTable(ctx context.Context, req CreateTableRequest) error
	AddColumn(ctx context.Context, table string, col ColumnSpec) error
	DropColumn(ctx context.Context, table string, colName string) error
	GenerateDDL(ctx context.Context, tableName string) (string, error)
	GenerateFullDDL(ctx context.Context) (string, error)
	Config() config.ConnectionConfig
}
