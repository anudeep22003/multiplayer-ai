package health

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubDB struct {
	err error
}

func (s stubDB) Ping(_ context.Context) error {
	return s.err
}

func TestNewMux(t *testing.T) {
	tests := []struct {
		name       string
		db         DB
		method     string
		target     string
		wantStatus int
	}{
		{
			name:       "GET health succeeds",
			method:     http.MethodGet,
			target:     "/health",
			wantStatus: http.StatusOK,
		},
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
			name:       "GET ready succeeds when db healthy",
			db:         stubDB{},
			method:     http.MethodGet,
			target:     "/ready",
			wantStatus: http.StatusOK,
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

			mux := NewMux(tc.db)
			req := httptest.NewRequest(tc.method, tc.target, nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatus)
			}
		})
	}
}
