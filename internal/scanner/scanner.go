package scanner

import (
	"context"
	"fmt"

	"db-studio-go/internal/config"
)

// Scanner defines the interface for detecting database configurations
type Scanner interface {
	Name() string
	Scan(ctx context.Context, projectPath string) ([]config.ConnectionConfig, error)
}

// CompositeScanner runs multiple scanners in sequence
type CompositeScanner struct {
	scanners []Scanner
}

// NewCompositeScanner creates a pipeline of scanners
func NewCompositeScanner(scanners ...Scanner) *CompositeScanner {
	return &CompositeScanner{scanners: scanners}
}

// Scan runs all registered scanners and collects detected database connections
func (cs *CompositeScanner) Scan(ctx context.Context, projectPath string) ([]config.ConnectionConfig, error) {
	var results []config.ConnectionConfig
	seenMap := make(map[string]bool)

	for _, s := range cs.scanners {
		conns, err := s.Scan(ctx, projectPath)
		if err != nil {
			// Continue scanning with other scanners if one fails
			continue
		}

		for _, conn := range conns {
			key := fmt.Sprintf("%s:%s:%s:%d:%s:%s", conn.Name, conn.Driver, conn.Host, conn.Port, conn.Database, conn.FilePath)
			if !seenMap[key] {
				seenMap[key] = true
				results = append(results, conn)
			}
		}
	}

	return results, nil
}
