package server

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/anudeep22003/conman/internal"
)

type stubDB struct {
	pgx.Tx // Supplies unused database methods for readiness tests.
	err    error
}

func (s stubDB) Ping(_ context.Context) error {
	return s.err
}

func TestNew(t *testing.T) {
	t.Parallel()

	cfg := internal.Config{
		Host:         "127.0.0.1",
		Port:         9090,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	tagHeader := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Middleware", "applied")
			next.ServeHTTP(w, r)
		})
	}

	srv := New(&cfg, stubDB{}, PassThrough, tagHeader)

	if srv.Addr != "127.0.0.1:9090" {
		t.Errorf("Addr = %q, want %q", srv.Addr, "127.0.0.1:9090")
	}
	if srv.ReadTimeout != cfg.ReadTimeout || srv.WriteTimeout != cfg.WriteTimeout || srv.IdleTimeout != cfg.IdleTimeout {
		t.Errorf("timeouts = (%v, %v, %v), want (%v, %v, %v)",
			srv.ReadTimeout, srv.WriteTimeout, srv.IdleTimeout,
			cfg.ReadTimeout, cfg.WriteTimeout, cfg.IdleTimeout,
		)
	}

	tests := []struct {
		name       string
		db         DB
		method     string
		target     string
		wantStatus int
	}{
		{
			name:       "GET healthz succeeds even when db fails",
			db:         stubDB{err: errors.New("db down")},
			method:     http.MethodGet,
			target:     "/healthz",
			wantStatus: http.StatusOK,
		},
		{
			name:       "POST healthz rejected",
			db:         stubDB{},
			method:     http.MethodPost,
			target:     "/healthz",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "GET readyz succeeds when db healthy",
			db:         stubDB{},
			method:     http.MethodGet,
			target:     "/readyz",
			wantStatus: http.StatusOK,
		},
		{
			name:       "GET readyz returns 503 when db ping fails",
			db:         stubDB{err: errors.New("connection refused")},
			method:     http.MethodGet,
			target:     "/readyz",
			wantStatus: http.StatusServiceUnavailable,
		},
		{
			name:       "GET readyz returns 503 when db is nil",
			db:         nil,
			method:     http.MethodGet,
			target:     "/readyz",
			wantStatus: http.StatusServiceUnavailable,
		},
		{
			name:       "unknown route returns 404",
			db:         stubDB{},
			method:     http.MethodGet,
			target:     "/unknown",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testSrv := New(&cfg, tc.db, PassThrough, tagHeader)
			req := httptest.NewRequest(tc.method, tc.target, nil)
			rec := httptest.NewRecorder()

			testSrv.Handler.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatus)
			}
			if got := rec.Header().Get("X-Middleware"); got != "applied" {
				t.Errorf("X-Middleware = %q, want %q", got, "applied")
			}
		})
	}
}
