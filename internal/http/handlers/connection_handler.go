package handlers

import (
	"net/http"

	"db-studio-go/internal/http/services"
)

type ConnectionHandler struct {
	BaseHandler
	service *services.ConnectionService
}

func NewConnectionHandler(service *services.ConnectionService) *ConnectionHandler {
	return &ConnectionHandler{service: service}
}

func (h *ConnectionHandler) HandleGetConnectionStatus(w http.ResponseWriter, r *http.Request) {
	statusResp, statusCode := h.service.GetStatus(r.Context())
	h.RespondJSON(w, statusCode, statusResp)
}
