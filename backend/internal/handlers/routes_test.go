package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5"
)

type stubDB struct {
	pgx.Tx // Supplies unused database methods for readiness tests.
	err    error
}

func (s stubDB) Ping(_ context.Context) error {
	return s.err
}

func TestNewRouter(t *testing.T) {
	tagHeader := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Middleware", "applied")
			next.ServeHTTP(w, r)
		})
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

			router := NewRouter(tc.db, tagHeader)
			req := httptest.NewRequest(tc.method, tc.target, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatus)
			}
			if got := rec.Header().Get("X-Middleware"); got != "applied" {
				t.Errorf("X-Middleware = %q, want %q", got, "applied")
			}
		})
	}
}

func TestLoginPreflight(t *testing.T) {
	router := NewRouter(nil, DevUserMiddleware)
	req := httptest.NewRequest(http.MethodOptions, "/login", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "content-type")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent || rec.Body.Len() != 0 {
		t.Fatalf("preflight returned %d with body %q, want 204 with empty body", rec.Code, rec.Body.String())
	}
}
