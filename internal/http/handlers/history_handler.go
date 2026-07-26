package handlers

import (
	"net/http"

	"db-studio-go/internal/http/services"
)

type HistoryHandler struct {
	BaseHandler
	service *services.HistoryService
}

func NewHistoryHandler(service *services.HistoryService) *HistoryHandler {
	return &HistoryHandler{service: service}
}

func (h *HistoryHandler) HandleGetHistory(w http.ResponseWriter, r *http.Request) {
	history, err := h.service.GetHistory()
	if err != nil {
		h.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"history": history,
	})
}

func (h *HistoryHandler) HandleClearHistory(w http.ResponseWriter, r *http.Request) {
	if err := h.service.ClearHistory(); err != nil {
		h.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "History cleared successfully",
	})
}
