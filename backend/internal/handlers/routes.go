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

// New registers the HTTP endpoints and applies middleware in the given order.
func New(db DB, middlewares ...Middleware) http.Handler {
	h := &Handler{db: db}
	mux := http.NewServeMux()
	mux.Handle("/login", h.authMux())
	mux.Handle("/", h.operationalMux())
	return Chain(middlewares...)(mux)
}
