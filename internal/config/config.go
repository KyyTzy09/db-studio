package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Manager handles reading and writing global OS DBStudio configuration
type Manager struct {
	configFilePath string
}

// NewManager initializes a Manager using OS UserConfigDir (~/.config/dbstudio/connections.json)
func NewManager() (*Manager, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		homeDir, hErr := os.UserHomeDir()
		if hErr != nil {
			return nil, fmt.Errorf("failed to locate config or home directory: %w", err)
		}
		configDir = filepath.Join(homeDir, ".config")
	}

	appConfigDir := filepath.Join(configDir, "dbstudio")
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory %s: %w", appConfigDir, err)
	}

	return &Manager{
		configFilePath: filepath.Join(appConfigDir, "connections.json"),
	}, nil
}

// GetConfigFilePath returns the absolute path to the connections.json file
func (m *Manager) GetConfigFilePath() string {
	return m.configFilePath
}

// Load reads and parses the connections.json file
func (m *Manager) Load() (*GlobalConfig, error) {
	if _, err := os.Stat(m.configFilePath); errors.Is(err, os.ErrNotExist) {
		return &GlobalConfig{Connections: []ConnectionConfig{}}, nil
	}

	data, err := os.ReadFile(m.configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if len(strings.TrimSpace(string(data))) == 0 {
		return &GlobalConfig{Connections: []ConnectionConfig{}}, nil
	}

	var cfg GlobalConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	return &cfg, nil
}

// Save writes the GlobalConfig to connections.json
func (m *Manager) Save(cfg *GlobalConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config JSON: %w", err)
	}

	if err := os.WriteFile(m.configFilePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// FindByProjectPath returns a saved connection matching projectPath (if exists)
func (m *Manager) FindByProjectPath(projectPath string) (*ConnectionConfig, bool, error) {
	cfg, err := m.Load()
	if err != nil {
		return nil, false, err
	}

	cleanTarget := filepath.Clean(projectPath)
	for _, conn := range cfg.Connections {
		if filepath.Clean(conn.ProjectPath) == cleanTarget {
			return &conn, true, nil
		}
	}

	return nil, false, nil
}

// SaveConnection adds or updates a connection for a project
func (m *Manager) SaveConnection(conn ConnectionConfig) error {
	cfg, err := m.Load()
	if err != nil {
		return err
	}

	now := time.Now()
	conn.UpdatedAt = now
	cleanTarget := filepath.Clean(conn.ProjectPath)

	foundIndex := -1
	for i, c := range cfg.Connections {
		if filepath.Clean(c.ProjectPath) == cleanTarget || (conn.ID != "" && c.ID == conn.ID) {
			foundIndex = i
			break
		}
	}

	if conn.ID == "" {
		conn.ID = fmt.Sprintf("conn_%d", now.UnixNano())
	}
	if conn.CreatedAt.IsZero() {
		conn.CreatedAt = now
	}

	if foundIndex >= 0 {
		cfg.Connections[foundIndex] = conn
	} else {
		cfg.Connections = append(cfg.Connections, conn)
	}

	cfg.ActiveConnectionID = conn.ID
	return m.Save(cfg)
}

// AddHistoryItem appends a history item to the specified connection and limits to 200 items
func (m *Manager) AddHistoryItem(connID string, item QueryHistoryItem) error {
	cfg, err := m.Load()
	if err != nil {
		return err
	}

	for i, conn := range cfg.Connections {
		if conn.ID == connID {
			// Prepend new item so newest is first
			history := append([]QueryHistoryItem{item}, conn.History...)
			if len(history) > 200 {
				history = history[:200]
			}
			cfg.Connections[i].History = history
			return m.Save(cfg)
		}
	}
	return nil
}

// GetHistory returns history items for a connection
func (m *Manager) GetHistory(connID string) ([]QueryHistoryItem, error) {
	cfg, err := m.Load()
	if err != nil {
		return nil, err
	}

	for _, conn := range cfg.Connections {
		if conn.ID == connID {
			if conn.History == nil {
				return []QueryHistoryItem{}, nil
			}
			return conn.History, nil
		}
	}
	return []QueryHistoryItem{}, nil
}

// ClearHistory clears all history items for a connection
func (m *Manager) ClearHistory(connID string) error {
	cfg, err := m.Load()
	if err != nil {
		return err
	}

	for i, conn := range cfg.Connections {
		if conn.ID == connID {
			cfg.Connections[i].History = []QueryHistoryItem{}
			return m.Save(cfg)
		}
	}
	return nil
}

// GetSnippets returns saved query snippets for a connection
func (m *Manager) GetSnippets(connID string) ([]QuerySnippet, error) {
	cfg, err := m.Load()
	if err != nil {
		return nil, err
	}

	for _, conn := range cfg.Connections {
		if conn.ID == connID {
			if conn.Snippets == nil {
				return []QuerySnippet{}, nil
			}
			return conn.Snippets, nil
		}
	}
	return []QuerySnippet{}, nil
}

// SaveSnippet adds or updates a query snippet for a connection
func (m *Manager) SaveSnippet(connID string, snippet QuerySnippet) error {
	cfg, err := m.Load()
	if err != nil {
		return err
	}

	now := time.Now()
	if snippet.ID == "" {
		snippet.ID = fmt.Sprintf("snip_%d", now.UnixNano())
		snippet.CreatedAt = now
	}
	snippet.UpdatedAt = now

	for i, conn := range cfg.Connections {
		if conn.ID == connID {
			foundIdx := -1
			for j, s := range conn.Snippets {
				if s.ID == snippet.ID {
					foundIdx = j
					break
				}
			}
			if foundIdx >= 0 {
				cfg.Connections[i].Snippets[foundIdx] = snippet
			} else {
				cfg.Connections[i].Snippets = append([]QuerySnippet{snippet}, conn.Snippets...)
			}
			return m.Save(cfg)
		}
	}
	return nil
}

// DeleteSnippet removes a snippet by ID for a connection
func (m *Manager) DeleteSnippet(connID string, snippetID string) error {
	cfg, err := m.Load()
	if err != nil {
		return err
	}

	for i, conn := range cfg.Connections {
		if conn.ID == connID {
			var newSnippets []QuerySnippet
			for _, s := range conn.Snippets {
				if s.ID != snippetID {
					newSnippets = append(newSnippets, s)
				}
			}
			cfg.Connections[i].Snippets = newSnippets
			return m.Save(cfg)
		}
	}
	return nil
}
