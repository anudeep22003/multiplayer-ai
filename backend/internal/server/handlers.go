package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

const readinessTimeout = 2 * time.Second

// DB defines the database operations required by the server's root handlers.
type DB interface {
	Ping(ctx context.Context) error
	QueryRow(context.Context, string, ...any) pgx.Row
}

// Handler holds dependencies for the server's HTTP endpoints.
type Handler struct {
	db DB
}

// NewHandler constructs a Handler with its required dependencies.
func NewHandler(db DB) *Handler {
	return &Handler{db: db}
}

// Healthz reports process liveness without checking external dependencies.
func (h *Handler) Healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Readyz verifies downstream dependencies (PostgreSQL) can serve traffic.
func (h *Handler) Readyz(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), readinessTimeout)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}

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
