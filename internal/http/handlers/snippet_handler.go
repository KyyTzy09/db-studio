package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"db-studio-go/internal/config"
	"db-studio-go/internal/http/services"
)

type SnippetHandler struct {
	BaseHandler
	service *services.SnippetService
}

func NewSnippetHandler(service *services.SnippetService) *SnippetHandler {
	return &SnippetHandler{service: service}
}

func (h *SnippetHandler) HandleGetSnippets(w http.ResponseWriter, r *http.Request) {
	snippets, err := h.service.GetSnippets()
	if err != nil {
		h.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"snippets": snippets,
	})
}

func (h *SnippetHandler) HandleSaveSnippet(w http.ResponseWriter, r *http.Request) {
	var snippet config.QuerySnippet
	if err := json.NewDecoder(r.Body).Decode(&snippet); err != nil {
		h.RespondError(w, http.StatusBadRequest, "Invalid snippet JSON payload")
		return
	}

	if snippet.Title == "" || snippet.Query == "" {
		h.RespondError(w, http.StatusBadRequest, "Snippet title and query are required")
		return
	}

	if err := h.service.SaveSnippet(snippet); err != nil {
		h.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Snippet saved successfully",
	})
}

func (h *SnippetHandler) HandleDeleteSnippet(w http.ResponseWriter, r *http.Request) {
	snippetID := chi.URLParam(r, "id")
	if snippetID == "" {
		h.RespondError(w, http.StatusBadRequest, "Snippet ID is required")
		return
	}

	if err := h.service.DeleteSnippet(snippetID); err != nil {
		h.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Snippet deleted successfully",
	})
}
