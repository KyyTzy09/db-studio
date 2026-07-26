package services

import (
	"db-studio-go/internal/config"
	"db-studio-go/internal/db"
)

type SnippetService struct {
	driver        db.Database
	configManager *config.Manager
}

func NewSnippetService(driver db.Database, configManager *config.Manager) *SnippetService {
	return &SnippetService{
		driver:        driver,
		configManager: configManager,
	}
}

func (s *SnippetService) GetSnippets() ([]config.QuerySnippet, error) {
	if s.configManager == nil {
		return []config.QuerySnippet{}, nil
	}
	connID := s.driver.Config().ID
	return s.configManager.GetSnippets(connID)
}

func (s *SnippetService) SaveSnippet(snippet config.QuerySnippet) error {
	if s.configManager == nil {
		return nil
	}
	connID := s.driver.Config().ID
	return s.configManager.SaveSnippet(connID, snippet)
}

func (s *SnippetService) DeleteSnippet(snippetID string) error {
	if s.configManager == nil {
		return nil
	}
	connID := s.driver.Config().ID
	return s.configManager.DeleteSnippet(connID, snippetID)
}
