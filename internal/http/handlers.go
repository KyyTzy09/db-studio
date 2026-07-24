package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"db-studio-go/internal/db"
)

type Handler struct {
	driver db.Database
}

func NewHandler(driver db.Database) *Handler {
	return &Handler{driver: driver}
}

func (h *Handler) respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *Handler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}

// GET /api/connection/status (Triggers lazy connection)
func (h *Handler) HandleGetConnectionStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := h.driver.Ping(ctx)
	if err != nil {
		h.respondJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
			"connected": false,
			"error":     err.Error(),
			"config":    h.driver.Config(),
		})
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"connected": true,
		"config":    h.driver.Config(),
	})
}

// GET /api/tables
func (h *Handler) HandleGetTables(w http.ResponseWriter, r *http.Request) {
	tables, err := h.driver.GetTables(r.Context())
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if tables == nil {
		tables = []db.TableInfo{}
	}
	h.respondJSON(w, http.StatusOK, map[string]interface{}{"tables": tables})
}

// GET /api/tables/{name}/schema
func (h *Handler) HandleGetSchema(w http.ResponseWriter, r *http.Request) {
	tableName := chi.URLParam(r, "name")
	schema, err := h.driver.GetSchema(r.Context(), tableName)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.respondJSON(w, http.StatusOK, schema)
}

// GET /api/tables/{name}/data
func (h *Handler) HandleGetData(w http.ResponseWriter, r *http.Request) {
	tableName := chi.URLParam(r, "name")
	queryStr := "SELECT * FROM " + tableName + " LIMIT 100;"

	result, err := h.driver.ExecuteQuery(r.Context(), queryStr, true)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.respondJSON(w, http.StatusOK, result)
}

// POST /api/tables/{name}
func (h *Handler) HandleInsertRow(w http.ResponseWriter, r *http.Request) {
	tableName := chi.URLParam(r, "name")
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if err := h.driver.InsertRow(r.Context(), tableName, data); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.respondJSON(w, http.StatusCreated, map[string]string{"message": "Row inserted successfully"})
}

// PATCH /api/tables/{name}
func (h *Handler) HandleUpdateRow(w http.ResponseWriter, r *http.Request) {
	tableName := chi.URLParam(r, "name")
	var payload struct {
		PK   map[string]interface{} `json:"pk"`
		Data map[string]interface{} `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if err := h.driver.UpdateRow(r.Context(), tableName, payload.PK, payload.Data); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Row updated successfully"})
}

// DELETE /api/tables/{name}
func (h *Handler) HandleDeleteRow(w http.ResponseWriter, r *http.Request) {
	tableName := chi.URLParam(r, "name")
	var pk map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&pk); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if err := h.driver.DeleteRow(r.Context(), tableName, pk); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Row deleted successfully"})
}

type QueryPayload struct {
	Query string `json:"query"`
	Force bool   `json:"force"`
}

// POST /api/query (Supports HTTP 428 Precondition Required Safety Guard)
func (h *Handler) HandleExecuteQuery(w http.ResponseWriter, r *http.Request) {
	var payload QueryPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	result, err := h.driver.ExecuteQuery(r.Context(), payload.Query, payload.Force)
	if err != nil {
		if strings.Contains(err.Error(), "DESTRUCTIVE_QUERY_WARNING") {
			h.respondJSON(w, http.StatusPreconditionRequired, map[string]interface{}{
				"warning": "Kueri ini mengandung kata kunci destruktif (DELETE/UPDATE/DROP/TRUNCATE). Lanjutkan?",
				"code":    "DESTRUCTIVE_QUERY",
			})
			return
		}
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}
