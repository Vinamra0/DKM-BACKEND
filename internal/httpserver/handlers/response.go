package handlers

import (
	"encoding/json"
	"net/http"
)

// APIResponse is the standard response envelope for all API endpoints
type APIResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

// WriteJSON writes a JSON response with the standard envelope
func WriteJSON(w http.ResponseWriter, status int, success bool, data any, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := APIResponse{
		Success: success,
		Data:    data,
		Message: message,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// WriteSuccess writes a success response
func WriteSuccess(w http.ResponseWriter, status int, data any) {
	WriteJSON(w, status, true, data, "")
}

// WriteError writes an error response
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, false, nil, message)
}
