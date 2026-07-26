package services

import (
	"fmt"
	"time"

	"db-studio-go/internal/config"
	"db-studio-go/internal/db"
)

type HistoryService struct {
	driver        db.Database
	configManager *config.Manager
}

func NewHistoryService(driver db.Database, configManager *config.Manager) *HistoryService {
	return &HistoryService{
		driver:        driver,
		configManager: configManager,
	}
}

func (s *HistoryService) LogQuery(queryStr string, durationMs int64, status string, rowsAffected int64, errMsg string) error {
	if s.configManager == nil {
		return nil
	}
	connID := s.driver.Config().ID
	if connID == "" {
		return nil
	}

	item := config.QueryHistoryItem{
		ID:           fmt.Sprintf("hist_%d", time.Now().UnixNano()),
		Query:        queryStr,
		ExecutedAt:   time.Now(),
		DurationMs:   durationMs,
		Status:       status,
		RowsAffected: rowsAffected,
		ErrorMessage: errMsg,
	}

	return s.configManager.AddHistoryItem(connID, item)
}

func (s *HistoryService) GetHistory() ([]config.QueryHistoryItem, error) {
	if s.configManager == nil {
		return []config.QueryHistoryItem{}, nil
	}
	connID := s.driver.Config().ID
	return s.configManager.GetHistory(connID)
}

func (s *HistoryService) ClearHistory() error {
	if s.configManager == nil {
		return nil
	}
	connID := s.driver.Config().ID
	return s.configManager.ClearHistory(connID)
}
