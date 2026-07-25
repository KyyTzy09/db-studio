package models

import "db-studio-go/internal/db"

type TablesResponse struct {
	Tables []db.TableInfo `json:"tables"`
}

type UpdateRowPayload struct {
	PK   map[string]interface{} `json:"pk"`
	Data map[string]interface{} `json:"data"`
}

type BatchPayload struct {
	Rows []map[string]interface{} `json:"rows"`
	Mode string                   `json:"mode"`
}

type BatchResponse struct {
	Success      bool  `json:"success"`
	AffectedRows int64 `json:"affected_rows"`
}
