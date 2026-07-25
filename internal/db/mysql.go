package db

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"db-studio-go/internal/config"
)

type MySQLDriver struct {
	cfg config.ConnectionConfig
	db  *sql.DB
	mu  sync.Mutex
}

func NewMySQLDriver(cfg config.ConnectionConfig) *MySQLDriver {
	return &MySQLDriver{cfg: cfg}
}

func (m *MySQLDriver) Config() config.ConnectionConfig {
	return m.cfg
}

func (m *MySQLDriver) Connect(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.db != nil {
		return nil
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		m.cfg.User, m.cfg.Password, m.cfg.Host, m.cfg.Port, m.cfg.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open mysql connection: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return fmt.Errorf("mysql ping failed: %w", err)
	}

	m.db = db
	return nil
}

func (m *MySQLDriver) Disconnect() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.db != nil {
		err := m.db.Close()
		m.db = nil
		return err
	}
	return nil
}

func (m *MySQLDriver) Ping(ctx context.Context) error {
	if err := m.Connect(ctx); err != nil {
		return err
	}
	return m.db.PingContext(ctx)
}

func (m *MySQLDriver) GetTables(ctx context.Context) ([]TableInfo, error) {
	if err := m.Connect(ctx); err != nil {
		return nil, err
	}

	query := "SHOW FULL TABLES;"
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []TableInfo
	for rows.Next() {
		var name, tableType string
		if err := rows.Scan(&name, &tableType); err != nil {
			return nil, err
		}
		tables = append(tables, TableInfo{Name: name, Type: tableType})
	}

	return tables, nil
}

func (m *MySQLDriver) GetSchema(ctx context.Context, tableName string) (*TableSchema, error) {
	if err := m.Connect(ctx); err != nil {
		return nil, err
	}

	query := fmt.Sprintf("DESCRIBE `%s`;", tableName)
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var colName, dataType, isNull, keyStr string
		var defaultVal, extra sql.NullString

		if err := rows.Scan(&colName, &dataType, &isNull, &keyStr, &defaultVal, &extra); err != nil {
			return nil, err
		}

		columns = append(columns, ColumnInfo{
			Name:         colName,
			DataType:     dataType,
			IsNullable:   isNull == "YES",
			IsPrimaryKey: keyStr == "PRI",
			DefaultValue: defaultVal.String,
		})
	}

	return &TableSchema{TableName: tableName, Columns: columns}, nil
}

func (m *MySQLDriver) GetSchemaGraph(ctx context.Context) (*SchemaGraph, error) {
	tables, err := m.GetTables(ctx)
	if err != nil {
		return nil, err
	}

	var nodes []TableSchema
	for _, t := range tables {
		if t.Type != "BASE TABLE" {
			continue
		}
		schema, err := m.GetSchema(ctx, t.Name)
		if err != nil {
			continue
		}
		nodes = append(nodes, *schema)
	}

	query := `
		SELECT 
			TABLE_NAME AS source_table,
			COLUMN_NAME AS source_column,
			REFERENCED_TABLE_NAME AS target_table,
			REFERENCED_COLUMN_NAME AS target_column
		FROM information_schema.KEY_COLUMN_USAGE
		WHERE TABLE_SCHEMA = DATABASE() AND REFERENCED_TABLE_NAME IS NOT NULL;
	`

	rows, err := m.db.QueryContext(ctx, query)
	var edges []ForeignKeyRelation
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var rel ForeignKeyRelation
			if err := rows.Scan(&rel.SourceTable, &rel.SourceColumn, &rel.TargetTable, &rel.TargetColumn); err == nil {
				rel.ID = fmt.Sprintf("%s_%s-%s_%s", rel.SourceTable, rel.SourceColumn, rel.TargetTable, rel.TargetColumn)
				edges = append(edges, rel)
			}
		}
	}

	if nodes == nil {
		nodes = []TableSchema{}
	}
	if edges == nil {
		edges = []ForeignKeyRelation{}
	}

	return &SchemaGraph{Nodes: nodes, Edges: edges}, nil
}

func (m *MySQLDriver) ExecuteQuery(ctx context.Context, queryStr string, force bool) (*QueryResult, error) {
	if err := m.Connect(ctx); err != nil {
		return nil, err
	}

	destructivePattern := regexp.MustCompile(`(?i)\b(DROP|DELETE|UPDATE|TRUNCATE|ALTER)\b`)
	if !force && destructivePattern.MatchString(queryStr) {
		return nil, fmt.Errorf("DESTRUCTIVE_QUERY_WARNING: query contains dangerous keyword")
	}

	start := time.Now()
	rows, err := m.db.QueryContext(ctx, queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var resultRows []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{})
		for i, colName := range cols {
			val := values[i]
			if b, ok := val.([]byte); ok {
				rowMap[colName] = string(b)
			} else {
				rowMap[colName] = val
			}
		}
		resultRows = append(resultRows, rowMap)
	}

	if resultRows == nil {
		resultRows = []map[string]interface{}{}
	}

	elapsed := time.Since(start).Milliseconds()
	return &QueryResult{
		Columns:      cols,
		Rows:         resultRows,
		AffectedRows: int64(len(resultRows)),
		ExecutionMs:  elapsed,
	}, nil
}

func (m *MySQLDriver) InsertRow(ctx context.Context, table string, data map[string]interface{}) error {
	if err := m.Connect(ctx); err != nil {
		return err
	}

	if len(data) == 0 {
		query := fmt.Sprintf("INSERT INTO `%s` () VALUES ()", table)
		_, err := m.db.ExecContext(ctx, query)
		return err
	}

	var cols []string
	var placeholders []string
	var args []interface{}

	for k, v := range data {
		cols = append(cols, fmt.Sprintf("`%s`", k))
		placeholders = append(placeholders, "?")
		args = append(args, v)
	}

	query := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", table, strings.Join(cols, ", "), strings.Join(placeholders, ", "))
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}

func (m *MySQLDriver) UpdateRow(ctx context.Context, table string, pk map[string]interface{}, data map[string]interface{}) error {
	if err := m.Connect(ctx); err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	var setClauses []string
	var whereClauses []string
	var args []interface{}

	for k, v := range data {
		setClauses = append(setClauses, fmt.Sprintf("`%s` = ?", k))
		args = append(args, v)
	}

	for k, v := range pk {
		whereClauses = append(whereClauses, fmt.Sprintf("`%s` = ?", k))
		args = append(args, v)
	}

	query := fmt.Sprintf("UPDATE `%s` SET %s WHERE %s", table, strings.Join(setClauses, ", "), strings.Join(whereClauses, " AND "))
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}

func (m *MySQLDriver) DeleteRow(ctx context.Context, table string, pk map[string]interface{}) error {
	if err := m.Connect(ctx); err != nil {
		return err
	}

	var whereClauses []string
	var args []interface{}

	for k, v := range pk {
		whereClauses = append(whereClauses, fmt.Sprintf("`%s` = ?", k))
		args = append(args, v)
	}

	query := fmt.Sprintf("DELETE FROM `%s` WHERE %s", table, strings.Join(whereClauses, " AND "))
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}

func (m *MySQLDriver) BatchInsertOrUpdate(ctx context.Context, table string, rows []map[string]interface{}, mode string) (int64, error) {
	if len(rows) == 0 {
		return 0, nil
	}
	schema, err := m.GetSchema(ctx, table)
	if err != nil {
		return 0, err
	}

	var pkCols []string
	for _, col := range schema.Columns {
		if col.IsPrimaryKey {
			pkCols = append(pkCols, col.Name)
		}
	}

	var totalAffected int64
	for _, rowData := range rows {
		if mode == "upsert" && len(pkCols) > 0 {
			pkMap := make(map[string]interface{})
			hasPKValues := true
			for _, pkCol := range pkCols {
				val, exists := rowData[pkCol]
				if exists && val != nil && fmt.Sprintf("%v", val) != "" {
					pkMap[pkCol] = val
				} else {
					hasPKValues = false
				}
			}

			if hasPKValues {
				err := m.UpdateRow(ctx, table, pkMap, rowData)
				if err == nil {
					totalAffected++
					continue
				}
			}
		}

		err := m.InsertRow(ctx, table, rowData)
		if err != nil {
			return totalAffected, fmt.Errorf("row insert error: %w", err)
		}
		totalAffected++
	}

	return totalAffected, nil
}

func (m *MySQLDriver) CreateTable(ctx context.Context, req CreateTableRequest) error {
	if err := m.Connect(ctx); err != nil {
		return err
	}
	if req.TableName == "" || len(req.Columns) == 0 {
		return fmt.Errorf("table name and at least one column are required")
	}

	var colDefs []string
	var pks []string

	for _, c := range req.Columns {
		def := fmt.Sprintf("`%s` %s", c.Name, c.DataType)
		if c.AutoIncrement {
			def += " AUTO_INCREMENT"
		}
		if !c.IsNullable {
			def += " NOT NULL"
		}
		if c.DefaultValue != "" {
			def += fmt.Sprintf(" DEFAULT %s", c.DefaultValue)
		}
		colDefs = append(colDefs, def)
		if c.IsPrimaryKey {
			pks = append(pks, fmt.Sprintf("`%s`", c.Name))
		}
	}

	if len(pks) > 0 {
		colDefs = append(colDefs, fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(pks, ", ")))
	}

	query := fmt.Sprintf("CREATE TABLE `%s` (\n  %s\n)", req.TableName, strings.Join(colDefs, ",\n  "))
	_, err := m.db.ExecContext(ctx, query)
	return err
}

func (m *MySQLDriver) AddColumn(ctx context.Context, table string, col ColumnSpec) error {
	if err := m.Connect(ctx); err != nil {
		return err
	}
	def := fmt.Sprintf("`%s` %s", col.Name, col.DataType)
	if !col.IsNullable {
		def += " NOT NULL"
	}
	if col.DefaultValue != "" {
		def += fmt.Sprintf(" DEFAULT %s", col.DefaultValue)
	}

	query := fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s", table, def)
	_, err := m.db.ExecContext(ctx, query)
	return err
}

func (m *MySQLDriver) DropColumn(ctx context.Context, table string, colName string) error {
	if err := m.Connect(ctx); err != nil {
		return err
	}
	query := fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`", table, colName)
	_, err := m.db.ExecContext(ctx, query)
	return err
}
