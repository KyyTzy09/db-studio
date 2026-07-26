package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"db-studio-go/internal/http/services"
)

type DDLHandler struct {
	BaseHandler
	service *services.DDLService
}

func NewDDLHandler(service *services.DDLService) *DDLHandler {
	return &DDLHandler{service: service}
}

func (h *DDLHandler) HandleGetTableDDL(w http.ResponseWriter, r *http.Request) {
	tableName := chi.URLParam(r, "name")
	if tableName == "" {
		h.RespondError(w, http.StatusBadRequest, "Table name is required")
		return
	}

	ddl, err := h.service.GenerateTableDDL(r.Context(), tableName)
	if err != nil {
		h.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"table_name": tableName,
		"ddl":        ddl,
	})
}

func (h *DDLHandler) HandleExportFullDDL(w http.ResponseWriter, r *http.Request) {
	ddl, err := h.service.GenerateFullDDL(r.Context())
	if err != nil {
		h.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/sql")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"schema_dump.sql\""))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(ddl))
}
