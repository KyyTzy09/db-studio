package services

import (
	"context"

	"db-studio-go/internal/db"
	"db-studio-go/internal/http/models"
)

type QueryService struct {
	driver db.Database
}

func NewQueryService(driver db.Database) *QueryService {
	return &QueryService{driver: driver}
}

func (s *QueryService) ExecuteQuery(ctx context.Context, payload models.QueryPayload) (*db.QueryResult, error) {
	return s.driver.ExecuteQuery(ctx, payload.Query, payload.Force)
}
