package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// API Gateway supplies the CORS headers on the preflight response.
func loginPreflight(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

const upsertUser = `
 INSERT INTO users (username) VALUES ($1)
 ON CONFLICT (username) DO UPDATE SET username = EXCLUDED.username
 RETURNING id, username, (xmax = 0) AS is_new`

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil || decoder.Decode(new(any)) != io.EOF {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	username := strings.TrimSpace(input.Username)
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}
	var result struct {
		UserID   string `json:"userId"`
		Username string `json:"username"`
		IsNew    bool   `json:"isNew"`
	}
	if h.db == nil {
		http.Error(w, "database unavailable", http.StatusServiceUnavailable)
		return
	}
	err := h.db.QueryRow(r.Context(), upsertUser, username).Scan(&result.UserID, &result.Username, &result.IsNew)
	if err != nil {
		http.Error(w, "login failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
