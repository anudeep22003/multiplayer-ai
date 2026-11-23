package server

import "net/http"

// RegisterRoutes attaches the handler's endpoints to mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /login", h.Login)
	mux.HandleFunc("OPTIONS /login", loginPreflight)
	mux.HandleFunc("GET /health", h.Healthz)
	mux.HandleFunc("GET /healthz", h.Healthz)
	mux.HandleFunc("GET /readyz", h.Readyz)
}
