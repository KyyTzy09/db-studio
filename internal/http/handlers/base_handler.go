package handlers

import (
	"encoding/json"
	"net/http"

	"db-studio-go/internal/http/models"
)

type BaseHandler struct{}

func (h *BaseHandler) RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *BaseHandler) RespondError(w http.ResponseWriter, status int, message string) {
	h.RespondJSON(w, status, models.ErrorResponse{Error: message})
}
