package services

import (
	"context"
	"net/http"

	"db-studio-go/internal/db"
	"db-studio-go/internal/http/models"
)

type ConnectionService struct {
	driver db.Database
}

func NewConnectionService(driver db.Database) *ConnectionService {
	return &ConnectionService{driver: driver}
}

func (s *ConnectionService) GetStatus(ctx context.Context) (*models.ConnectionStatusResponse, int) {
	err := s.driver.Ping(ctx)
	if err != nil {
		return &models.ConnectionStatusResponse{
			Connected: false,
			Error:     err.Error(),
			Config:    s.driver.Config(),
		}, http.StatusServiceUnavailable
	}

	return &models.ConnectionStatusResponse{
		Connected: true,
		Config:    s.driver.Config(),
	}, http.StatusOK
}
