package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"db-studio-go/internal/db"
	"db-studio-go/internal/http/models"
	"db-studio-go/internal/http/services"
)

type TableHandler struct {
	BaseHandler
	service *services.TableService
}

func NewTableHandler(service *services.TableService) *TableHandler {
	return &TableHandler{service: service}
}

func (h *TableHandler) HandleGetTables(w http.ResponseWriter, r *http.Request) {
	tables, err := h.service.GetTables(r.Context())
	if err != nil {
		h.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.RespondJSON(w, http.StatusOK, models.TablesResponse{Tables: tables})
}

func (h *TableHandler) HandleGetSchema(w http.ResponseWriter, r *http.Request) {
	tableName := chi.URLParam(r, "name")
	schema, err := h.service.GetSchema(r.Context(), tableName)
	if err != nil {
		h.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.RespondJSON(w, http.StatusOK, schema)
}

func (h *TableHandler) HandleGetData(w http.ResponseWriter, r *http.Request) {
	tableName := chi.URLParam(r, "name")
	result, err := h.service.GetData(r.Context(), tableName)
	if err != nil {
		h.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.RespondJSON(w, http.StatusOK, result)
}

func (h *TableHandler) HandleInsertRow(w http.ResponseWriter, r *http.Request) {
	tableName := chi.URLParam(r, "name")
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.RespondError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if err := h.service.InsertRow(r.Context(), tableName, data); err != nil {
		h.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.RespondJSON(w, http.StatusCreated, models.MessageResponse{Message: "Row inserted successfully"})
}

func (h *TableHandler) HandleUpdateRow(w http.ResponseWriter, r *http.Request) {
	tableName := chi.URLParam(r, "name")
	var payload models.UpdateRowPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.RespondError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if err := h.service.UpdateRow(r.Context(), tableName, payload.PK, payload.Data); err != nil {
		h.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.RespondJSON(w, http.StatusOK, models.MessageResponse{Message: "Row updated successfully"})
}

func (h *TableHandler) HandleDeleteRow(w http.ResponseWriter, r *http.Request) {
	tableName := chi.URLParam(r, "name")
	var pk map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&pk); err != nil {
		h.RespondError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if err := h.service.DeleteRow(r.Context(), tableName, pk); err != nil {
		h.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.RespondJSON(w, http.StatusOK, models.MessageResponse{Message: "Row deleted successfully"})
}

func (h *TableHandler) HandleBatchInsertOrUpdate(w http.ResponseWriter, r *http.Request) {
	tableName := chi.URLParam(r, "name")
	var payload models.BatchPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.RespondError(w, http.StatusBadRequest, "Invalid JSON payload: "+err.Error())
		return
	}

	affected, err := h.service.BatchInsertOrUpdate(r.Context(), tableName, payload)
	if err != nil {
		h.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.RespondJSON(w, http.StatusOK, models.BatchResponse{
		Success:      true,
		AffectedRows: affected,
	})
}

func (h *TableHandler) HandleCreateTable(w http.ResponseWriter, r *http.Request) {
	var req db.CreateTableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.RespondError(w, http.StatusBadRequest, "Invalid JSON payload: "+err.Error())
		return
	}

	if err := h.service.CreateTable(r.Context(), req); err != nil {
		h.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.RespondJSON(w, http.StatusCreated, models.MessageResponse{Message: "Table created successfully"})
}

func (h *TableHandler) HandleAddColumn(w http.ResponseWriter, r *http.Request) {
	tableName := chi.URLParam(r, "name")
	var col db.ColumnSpec
	if err := json.NewDecoder(r.Body).Decode(&col); err != nil {
		h.RespondError(w, http.StatusBadRequest, "Invalid JSON payload: "+err.Error())
		return
	}

	if err := h.service.AddColumn(r.Context(), tableName, col); err != nil {
		h.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.RespondJSON(w, http.StatusOK, models.MessageResponse{Message: "Column added successfully"})
}

func (h *TableHandler) HandleDropColumn(w http.ResponseWriter, r *http.Request) {
	tableName := chi.URLParam(r, "name")
	colName := chi.URLParam(r, "col")

	if err := h.service.DropColumn(r.Context(), tableName, colName); err != nil {
		h.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.RespondJSON(w, http.StatusOK, models.MessageResponse{Message: "Column dropped successfully"})
}
