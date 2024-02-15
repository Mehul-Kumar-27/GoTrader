package routes

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	Success bool        `json:"success,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (m *Manager) writeJsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (m *Manager) readJsonRequest(r *http.Request, data interface{}) error {
	return json.NewDecoder(r.Body).Decode(data)
}
