package handlers

import (
	"encoding/json"
	"net/http"

	"dkmbackend/internal/models"
	"dkmbackend/internal/services"

	"github.com/go-chi/chi/v5"
)

type ProductHandler struct{ svc *services.ProductService }

func NewProductHandler(s *services.ProductService) *ProductHandler { return &ProductHandler{svc: s} }

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteSuccess(w, http.StatusOK, items)
}

func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	item, err := h.svc.Get(r.Context(), id)
	if err != nil || item == nil {
		WriteError(w, http.StatusNotFound, "product not found")
		return
	}
	WriteSuccess(w, http.StatusOK, item)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := h.svc.Create(r.Context(), &p); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	WriteSuccess(w, http.StatusCreated, p)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := h.svc.Update(r.Context(), id, &p); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	WriteSuccess(w, http.StatusOK, p)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
