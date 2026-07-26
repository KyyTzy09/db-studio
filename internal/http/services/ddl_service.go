package services

import (
	"context"

	"db-studio-go/internal/db"
)

type DDLService struct {
	driver db.Database
}

func NewDDLService(driver db.Database) *DDLService {
	return &DDLService{driver: driver}
}

func (s *DDLService) GenerateTableDDL(ctx context.Context, tableName string) (string, error) {
	return s.driver.GenerateDDL(ctx, tableName)
}

func (s *DDLService) GenerateFullDDL(ctx context.Context) (string, error) {
	return s.driver.GenerateFullDDL(ctx)
}
