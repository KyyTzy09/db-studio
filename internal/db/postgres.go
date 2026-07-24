package db

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"db-studio-go/internal/config"
)

type PostgresDriver struct {
	cfg config.ConnectionConfig
	db  *sql.DB
	mu  sync.Mutex
}

func NewPostgresDriver(cfg config.ConnectionConfig) *PostgresDriver {
	return &PostgresDriver{cfg: cfg}
}

func (p *PostgresDriver) Config() config.ConnectionConfig {
	return p.cfg
}

func (p *PostgresDriver) Connect(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.db != nil {
		return nil
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		p.cfg.User, p.cfg.Password, p.cfg.Host, p.cfg.Port, p.cfg.Database)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to open postgres connection: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return fmt.Errorf("postgres ping failed: %w", err)
	}

	p.db = db
	return nil
}

func (p *PostgresDriver) Disconnect() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.db != nil {
		err := p.db.Close()
		p.db = nil
		return err
	}
	return nil
}

func (p *PostgresDriver) Ping(ctx context.Context) error {
	if err := p.Connect(ctx); err != nil {
		return err
	}
	return p.db.PingContext(ctx)
}

func (p *PostgresDriver) GetTables(ctx context.Context) ([]TableInfo, error) {
	if err := p.Connect(ctx); err != nil {
		return nil, err
	}

	query := `
		SELECT table_name, table_type 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		ORDER BY table_name;
	`

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []TableInfo
	for rows.Next() {
		var t TableInfo
		if err := rows.Scan(&t.Name, &t.Type); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}

	return tables, nil
}

func (p *PostgresDriver) GetSchema(ctx context.Context, tableName string) (*TableSchema, error) {
	if err := p.Connect(ctx); err != nil {
		return nil, err
	}

	query := `
		SELECT 
			c.column_name, 
			c.data_type, 
			c.is_nullable = 'YES' as is_nullable,
			COALESCE(c.column_default, '') as default_value,
			CASE WHEN pk.column_name IS NOT NULL THEN true ELSE false END as is_primary_key
		FROM information_schema.columns c
		LEFT JOIN (
			SELECT kcu.column_name
			FROM information_schema.table_constraints tc
			JOIN information_schema.key_column_usage kcu 
			  ON tc.constraint_name = kcu.constraint_name
			WHERE tc.constraint_type = 'PRIMARY KEY' AND tc.table_name = $1
		) pk ON c.column_name = pk.column_name
		WHERE c.table_name = $1 AND c.table_schema = 'public'
		ORDER BY c.ordinal_position;
	`

	rows, err := p.db.QueryContext(ctx, query, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var col ColumnInfo
		if err := rows.Scan(&col.Name, &col.DataType, &col.IsNullable, &col.DefaultValue, &col.IsPrimaryKey); err != nil {
			return nil, err
		}
		columns = append(columns, col)
	}

	return &TableSchema{TableName: tableName, Columns: columns}, nil
}

func (p *PostgresDriver) ExecuteQuery(ctx context.Context, queryStr string, force bool) (*QueryResult, error) {
	if err := p.Connect(ctx); err != nil {
		return nil, err
	}

	// Safety check
	destructivePattern := regexp.MustCompile(`(?i)\b(DROP|DELETE|UPDATE|TRUNCATE|ALTER)\b`)
	if !force && destructivePattern.MatchString(queryStr) {
		return nil, fmt.Errorf("DESTRUCTIVE_QUERY_WARNING: query contains dangerous keyword")
	}

	start := time.Now()
	rows, err := p.db.QueryContext(ctx, queryStr)
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

func (p *PostgresDriver) InsertRow(ctx context.Context, table string, data map[string]interface{}) error {
	if err := p.Connect(ctx); err != nil {
		return err
	}

	if len(data) == 0 {
		query := fmt.Sprintf("INSERT INTO %s DEFAULT VALUES", table)
		_, err := p.db.ExecContext(ctx, query)
		return err
	}

	var cols []string
	var placeholders []string
	var args []interface{}
	i := 1

	for k, v := range data {
		cols = append(cols, k)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		args = append(args, v)
		i++
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(cols, ", "), strings.Join(placeholders, ", "))
	_, err := p.db.ExecContext(ctx, query, args...)
	return err
}

func (p *PostgresDriver) UpdateRow(ctx context.Context, table string, pk map[string]interface{}, data map[string]interface{}) error {
	if err := p.Connect(ctx); err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	var setClauses []string
	var whereClauses []string
	var args []interface{}
	i := 1

	for k, v := range data {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", k, i))
		args = append(args, v)
		i++
	}

	for k, v := range pk {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", k, i))
		args = append(args, v)
		i++
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, strings.Join(setClauses, ", "), strings.Join(whereClauses, " AND "))
	_, err := p.db.ExecContext(ctx, query, args...)
	return err
}

func (p *PostgresDriver) DeleteRow(ctx context.Context, table string, pk map[string]interface{}) error {
	if err := p.Connect(ctx); err != nil {
		return err
	}

	var whereClauses []string
	var args []interface{}
	i := 1

	for k, v := range pk {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", k, i))
		args = append(args, v)
		i++
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s", table, strings.Join(whereClauses, " AND "))
	_, err := p.db.ExecContext(ctx, query, args...)
	return err
}

func (p *PostgresDriver) BatchInsertOrUpdate(ctx context.Context, table string, rows []map[string]interface{}, mode string) (int64, error) {
	if len(rows) == 0 {
		return 0, nil
	}
	schema, err := p.GetSchema(ctx, table)
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
				err := p.UpdateRow(ctx, table, pkMap, rowData)
				if err == nil {
					totalAffected++
					continue
				}
			}
		}

		err := p.InsertRow(ctx, table, rowData)
		if err != nil {
			return totalAffected, fmt.Errorf("row insert error: %w", err)
		}
		totalAffected++
	}

	return totalAffected, nil
}
