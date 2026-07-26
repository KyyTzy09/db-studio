package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"db-studio-go/internal/http/models"
	"db-studio-go/internal/http/services"
)

type QueryHandler struct {
	BaseHandler
	service        *services.QueryService
	historyService *services.HistoryService
}

func NewQueryHandler(service *services.QueryService, historyService *services.HistoryService) *QueryHandler {
	return &QueryHandler{
		service:        service,
		historyService: historyService,
	}
}

func (h *QueryHandler) HandleExecuteQuery(w http.ResponseWriter, r *http.Request) {
	var payload models.QueryPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.RespondError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	result, err := h.service.ExecuteQuery(r.Context(), payload)
	if err != nil {
		if strings.Contains(err.Error(), "DESTRUCTIVE_QUERY_WARNING") {
			h.RespondJSON(w, http.StatusPreconditionRequired, models.DestructiveWarningResponse{
				Warning: "Kueri ini mengandung kata kunci destruktif (DELETE/UPDATE/DROP/TRUNCATE). Lanjutkan?",
				Code:    "DESTRUCTIVE_QUERY",
			})
			return
		}
		if h.historyService != nil {
			_ = h.historyService.LogQuery(payload.Query, 0, "error", 0, err.Error())
		}
		h.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if h.historyService != nil && result != nil {
		_ = h.historyService.LogQuery(payload.Query, result.ExecutionMs, "success", result.AffectedRows, "")
	}

	h.RespondJSON(w, http.StatusOK, result)
}
