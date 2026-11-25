package app

import (
	"net/http"

	"github.com/anudeep22003/multiplayer-ai/internal/db"
	"github.com/anudeep22003/multiplayer-ai/internal/handlers"
	"github.com/anudeep22003/multiplayer-ai/internal/handlers/auth"
	"github.com/anudeep22003/multiplayer-ai/internal/handlers/health"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewHandler wires all HTTP route groups onto a single root handler.
func NewHandler(pool *pgxpool.Pool) http.Handler {
	queries := db.New(pool)

	return handlers.Assemble(
		handlers.RouteGroup{
			Prefix: "", Handler: health.NewMux(pool),
		},
		handlers.RouteGroup{
			Prefix: "/auth", Handler: handlers.DevUserMiddleware(auth.NewMux(queries)),
		},
	)
}
