package handlers

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

// NewRouter registers the HTTP endpoints and applies middleware in the given order.
func NewRouter(db DB, middlewares ...Middleware) http.Handler {
	h := &Handler{db: db}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", h.Login)
	mux.HandleFunc("OPTIONS /login", loginPreflight)
	mux.HandleFunc("GET /health", h.Healthz)
	mux.HandleFunc("GET /healthz", h.Healthz)
	mux.HandleFunc("GET /readyz", h.Readyz)
	return Chain(middlewares...)(mux)
}
