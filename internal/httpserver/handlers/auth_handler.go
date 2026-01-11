package handlers

import (
	"encoding/json"
	"net/http"

	"dkmbackend/internal/services"
)

type AuthHandler struct{ svc *services.AuthService }

func NewAuthHandler(s *services.AuthService) *AuthHandler { return &AuthHandler{svc: s} }

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct{ Email, Password string }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	token, user, err := h.svc.Login(body.Email, body.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"token": token, "user": user})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth == "" || len(auth) < 8 { // "Bearer "
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	token := auth[7:]
	user, err := h.svc.ParseToken(token)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": user})
}

// small local helper to avoid import cycles
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
