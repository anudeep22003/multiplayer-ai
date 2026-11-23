package handlers

import (
	"context"
	"net/http"
	"time"
)

const readinessTimeout = 2 * time.Second

func (h *Handler) operationalMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Healthz)
	mux.HandleFunc("GET /healthz", h.Healthz)
	mux.HandleFunc("GET /readyz", h.Readyz)
	return mux
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
