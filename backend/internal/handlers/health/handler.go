package health

import (
	"context"
	"net/http"
	"time"
)

const readinessTimeout = time.Second * 5

type DB interface {
	Ping(ctx context.Context) error
}

type handler struct {
	db DB
}

func NewSrvMux(db DB) *http.ServeMux {
	h := &handler{db: db}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Healthz)
	mux.HandleFunc("GET /healthz", h.Healthz)
	mux.HandleFunc("GET /ready", h.Readyz)
	return mux
}

// handler methods
func (h *handler) Healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *handler) Readyz(w http.ResponseWriter, r *http.Request) {
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
