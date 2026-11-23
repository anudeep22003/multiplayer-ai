package server

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"
)

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
	mux.HandleFunc("OPTIONS /login", loginPreflight)
	mux.HandleFunc("GET /health", h.Healthz)
	mux.HandleFunc("GET /healthz", h.Healthz)
	mux.HandleFunc("GET /readyz", h.Readyz)
}
