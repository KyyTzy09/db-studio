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
	service *services.QueryService
}

func NewQueryHandler(service *services.QueryService) *QueryHandler {
	return &QueryHandler{service: service}
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
		h.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.RespondJSON(w, http.StatusOK, result)
}
