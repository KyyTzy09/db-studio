package services

import (
	"context"

	"db-studio-go/internal/db"
	"db-studio-go/internal/http/models"
)

type TableService struct {
	driver db.Database
}

func NewTableService(driver db.Database) *TableService {
	return &TableService{driver: driver}
}

func (s *TableService) GetTables(ctx context.Context) ([]db.TableInfo, error) {
	tables, err := s.driver.GetTables(ctx)
	if err != nil {
		return nil, err
	}
	if tables == nil {
		tables = []db.TableInfo{}
	}
	return tables, nil
}

func (s *TableService) GetSchema(ctx context.Context, tableName string) (*db.TableSchema, error) {
	return s.driver.GetSchema(ctx, tableName)
}

func (s *TableService) GetData(ctx context.Context, tableName string) (*db.QueryResult, error) {
	queryStr := "SELECT * FROM " + tableName + " LIMIT 100;"
	return s.driver.ExecuteQuery(ctx, queryStr, true)
}

func (s *TableService) InsertRow(ctx context.Context, tableName string, data map[string]interface{}) error {
	return s.driver.InsertRow(ctx, tableName, data)
}

func (s *TableService) UpdateRow(ctx context.Context, tableName string, pk, data map[string]interface{}) error {
	return s.driver.UpdateRow(ctx, tableName, pk, data)
}

func (s *TableService) DeleteRow(ctx context.Context, tableName string, pk map[string]interface{}) error {
	return s.driver.DeleteRow(ctx, tableName, pk)
}

func (s *TableService) BatchInsertOrUpdate(ctx context.Context, tableName string, payload models.BatchPayload) (int64, error) {
	mode := payload.Mode
	if mode == "" {
		mode = "insert"
	}
	return s.driver.BatchInsertOrUpdate(ctx, tableName, payload.Rows, mode)
}

func (s *TableService) CreateTable(ctx context.Context, req db.CreateTableRequest) error {
	return s.driver.CreateTable(ctx, req)
}

func (s *TableService) AddColumn(ctx context.Context, tableName string, col db.ColumnSpec) error {
	return s.driver.AddColumn(ctx, tableName, col)
}

func (s *TableService) DropColumn(ctx context.Context, tableName, colName string) error {
	return s.driver.DropColumn(ctx, tableName, colName)
}
