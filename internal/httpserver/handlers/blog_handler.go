package handlers

import (
	"encoding/json"
	"net/http"

	"dkmbackend/internal/models"
	"dkmbackend/internal/services"

	"github.com/go-chi/chi/v5"
)

type BlogHandler struct{ svc *services.BlogService }

func NewBlogHandler(s *services.BlogService) *BlogHandler { return &BlogHandler{svc: s} }

func (h *BlogHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": items})
}

func (h *BlogHandler) Get(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "param")
	item, err := h.svc.GetByParam(r.Context(), param)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": item})
}

func (h *BlogHandler) Create(w http.ResponseWriter, r *http.Request) {
	var b models.Blog
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}
	if err := h.svc.Create(r.Context(), &b); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"blog": b})
}

func (h *BlogHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var b models.Blog
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}
	if err := h.svc.Update(r.Context(), id, &b); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"blog": b})
}

func (h *BlogHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
