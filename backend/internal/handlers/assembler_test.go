package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anudeep22003/conman/internal/handlers"
	"github.com/anudeep22003/conman/internal/handlers/auth"
	"github.com/anudeep22003/conman/internal/handlers/health"
)

func TestAssemble(t *testing.T) {
	pingMux := http.NewServeMux()
	pingMux.HandleFunc("GET /ping", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux := handlers.Assemble(
		handlers.RouteGroup{Prefix: "", Handler: health.NewMux(nil)},
		handlers.RouteGroup{Prefix: "/auth", Handler: auth.NewMux(nil)},
		handlers.RouteGroup{Prefix: "v1/", Handler: pingMux},
	)

	tests := []struct {
		name       string
		method     string
		target     string
		wantStatus int
	}{
		{
			name:       "root prefix routes to health mux",
			method:     http.MethodGet,
			target:     "/health",
			wantStatus: http.StatusOK,
		},
		{
			name:       "auth prefix routes to auth mux",
			method:     http.MethodOptions,
			target:     "/auth/login",
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "auth prefix enforces method from sub-mux",
			method:     http.MethodGet,
			target:     "/auth/login",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "prefix without leading slash and with trailing slash is normalized",
			method:     http.MethodGet,
			target:     "/v1/ping",
			wantStatus: http.StatusOK,
		},
		{
			name:       "unknown route returns 404",
			method:     http.MethodGet,
			target:     "/unknown",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tc.method, tc.target, nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatus)
			}
		})
	}
}
