package db

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"

	"db-studio-go/internal/config"
)

type SQLiteDriver struct {
	cfg config.ConnectionConfig
	db  *sql.DB
	mu  sync.Mutex
}

func NewSQLiteDriver(cfg config.ConnectionConfig) *SQLiteDriver {
	return &SQLiteDriver{cfg: cfg}
}

func (s *SQLiteDriver) Config() config.ConnectionConfig {
	return s.cfg
}

func (s *SQLiteDriver) Connect(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db != nil {
		return nil
	}

	db, err := sql.Open("sqlite", s.cfg.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open sqlite connection: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return fmt.Errorf("sqlite ping failed: %w", err)
	}

	s.db = db
	return nil
}

func (s *SQLiteDriver) Disconnect() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db != nil {
		err := s.db.Close()
		s.db = nil
		return err
	}
	return nil
}

func (s *SQLiteDriver) Ping(ctx context.Context) error {
	if err := s.Connect(ctx); err != nil {
		return err
	}
	return s.db.PingContext(ctx)
}

func (s *SQLiteDriver) GetTables(ctx context.Context) ([]TableInfo, error) {
	if err := s.Connect(ctx); err != nil {
		return nil, err
	}

	query := "SELECT name, type FROM sqlite_master WHERE type IN ('table', 'view') AND name NOT LIKE 'sqlite_%' ORDER BY name;"
	rows, err := s.db.QueryContext(ctx, query)
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
		tables = append(tables, TableInfo{Name: name, Type: strings.ToUpper(tableType)})
	}

	return tables, nil
}

func (s *SQLiteDriver) GetSchema(ctx context.Context, tableName string) (*TableSchema, error) {
	if err := s.Connect(ctx); err != nil {
		return nil, err
	}

	query := fmt.Sprintf("PRAGMA table_info(`%s`);", tableName)
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var cid int
		var name, dataType string
		var notNull, pk int
		var dfltValue sql.NullString

		if err := rows.Scan(&cid, &name, &dataType, &notNull, &dfltValue, &pk); err != nil {
			return nil, err
		}

		columns = append(columns, ColumnInfo{
			Name:         name,
			DataType:     dataType,
			IsNullable:   notNull == 0,
			IsPrimaryKey: pk > 0,
			DefaultValue: dfltValue.String,
		})
	}

	return &TableSchema{TableName: tableName, Columns: columns}, nil
}

func (s *SQLiteDriver) ExecuteQuery(ctx context.Context, queryStr string, force bool) (*QueryResult, error) {
	if err := s.Connect(ctx); err != nil {
		return nil, err
	}

	destructivePattern := regexp.MustCompile(`(?i)\b(DROP|DELETE|UPDATE|TRUNCATE|ALTER)\b`)
	if !force && destructivePattern.MatchString(queryStr) {
		return nil, fmt.Errorf("DESTRUCTIVE_QUERY_WARNING: query contains dangerous keyword")
	}

	start := time.Now()
	rows, err := s.db.QueryContext(ctx, queryStr)
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

func (s *SQLiteDriver) InsertRow(ctx context.Context, table string, data map[string]interface{}) error {
	if err := s.Connect(ctx); err != nil {
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
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *SQLiteDriver) UpdateRow(ctx context.Context, table string, pk map[string]interface{}, data map[string]interface{}) error {
	if err := s.Connect(ctx); err != nil {
		return err
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
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *SQLiteDriver) DeleteRow(ctx context.Context, table string, pk map[string]interface{}) error {
	if err := s.Connect(ctx); err != nil {
		return err
	}

	var whereClauses []string
	var args []interface{}

	for k, v := range pk {
		whereClauses = append(whereClauses, fmt.Sprintf("`%s` = ?", k))
		args = append(args, v)
	}

	query := fmt.Sprintf("DELETE FROM `%s` WHERE %s", table, strings.Join(whereClauses, " AND "))
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *SQLiteDriver) BatchInsertOrUpdate(ctx context.Context, table string, rows []map[string]interface{}, mode string) (int64, error) {
	if len(rows) == 0 {
		return 0, nil
	}
	schema, err := s.GetSchema(ctx, table)
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
				err := s.UpdateRow(ctx, table, pkMap, rowData)
				if err == nil {
					totalAffected++
					continue
				}
			}
		}

		err := s.InsertRow(ctx, table, rowData)
		if err != nil {
			return totalAffected, fmt.Errorf("row insert error: %w", err)
		}
		totalAffected++
	}

	return totalAffected, nil
}

func (s *SQLiteDriver) CreateTable(ctx context.Context, req CreateTableRequest) error {
	if err := s.Connect(ctx); err != nil {
		return err
	}
	if req.TableName == "" || len(req.Columns) == 0 {
		return fmt.Errorf("table name and at least one column are required")
	}

	var colDefs []string
	var pks []string

	for _, c := range req.Columns {
		def := fmt.Sprintf("\"%s\" %s", c.Name, c.DataType)
		if c.IsPrimaryKey && c.AutoIncrement {
			def = fmt.Sprintf("\"%s\" INTEGER PRIMARY KEY AUTOINCREMENT", c.Name)
		} else {
			if !c.IsNullable {
				def += " NOT NULL"
			}
			if c.DefaultValue != "" {
				def += fmt.Sprintf(" DEFAULT %s", c.DefaultValue)
			}
			if c.IsPrimaryKey {
				pks = append(pks, fmt.Sprintf("\"%s\"", c.Name))
			}
		}
		colDefs = append(colDefs, def)
	}

	if len(pks) > 0 {
		colDefs = append(colDefs, fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(pks, ", ")))
	}

	query := fmt.Sprintf("CREATE TABLE \"%s\" (\n  %s\n)", req.TableName, strings.Join(colDefs, ",\n  "))
	_, err := s.db.ExecContext(ctx, query)
	return err
}

func (s *SQLiteDriver) AddColumn(ctx context.Context, table string, col ColumnSpec) error {
	if err := s.Connect(ctx); err != nil {
		return err
	}
	def := fmt.Sprintf("\"%s\" %s", col.Name, col.DataType)
	if !col.IsNullable {
		def += " NOT NULL"
	}
	if col.DefaultValue != "" {
		def += fmt.Sprintf(" DEFAULT %s", col.DefaultValue)
	}

	query := fmt.Sprintf("ALTER TABLE \"%s\" ADD COLUMN %s", table, def)
	_, err := s.db.ExecContext(ctx, query)
	return err
}

func (s *SQLiteDriver) DropColumn(ctx context.Context, table string, colName string) error {
	if err := s.Connect(ctx); err != nil {
		return err
	}
	query := fmt.Sprintf("ALTER TABLE \"%s\" DROP COLUMN \"%s\"", table, colName)
	_, err := s.db.ExecContext(ctx, query)
	return err
}
