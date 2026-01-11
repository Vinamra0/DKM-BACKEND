package handlers

import (
	"encoding/json"
	"net/http"

	"dkmbackend/internal/models"
	"dkmbackend/internal/services"

	"github.com/go-chi/chi/v5"
)

type CareerHandler struct{ svc *services.CareerService }

func NewCareerHandler(s *services.CareerService) *CareerHandler { return &CareerHandler{svc: s} }

func (h *CareerHandler) PublicList(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.Public(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": items})
}

func (h *CareerHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": items})
}

func (h *CareerHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	item, err := h.svc.Get(r.Context(), id)
	if err != nil || item == nil {
		http.Error(w, "not found", 404)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": item})
}

func (h *CareerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var c models.Career
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}
	if err := h.svc.Create(r.Context(), &c); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"career": c})
}

func (h *CareerHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var c models.Career
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}
	if err := h.svc.Update(r.Context(), id, &c); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"career": c})
}

func (h *CareerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
