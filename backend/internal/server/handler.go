package server

import (
	"context"
	"net/http"
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

// RegisterRoutes attaches the handler's endpoints to mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /login", h.Login)
	// API Gateway supplies the CORS headers on the preflight response.
	mux.HandleFunc("OPTIONS /login", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("GET /health", h.Healthz)
	mux.HandleFunc("GET /healthz", h.Healthz)
	mux.HandleFunc("GET /readyz", h.Readyz)
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
