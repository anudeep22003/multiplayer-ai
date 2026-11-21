package server

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestPassThrough(t *testing.T) {
	t.Parallel()

	called := false
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if got := r.Header.Get("X-Test"); got != "ping" {
			t.Errorf("request header X-Test = %q, want %q", got, "ping")
		}
		w.Header().Set("X-Reply", "pong")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("ok"))
	})

	req := httptest.NewRequest(http.MethodPost, "/items", nil)
	req.Header.Set("X-Test", "ping")
	rec := httptest.NewRecorder()

	PassThrough(base).ServeHTTP(rec, req)

	if !called {
		t.Fatal("expected next handler to be called")
	}
	if rec.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusCreated)
	}
	if got := rec.Header().Get("X-Reply"); got != "pong" {
		t.Errorf("response header X-Reply = %q, want %q", got, "pong")
	}
	if got := rec.Body.String(); got != "ok" {
		t.Errorf("body = %q, want %q", got, "ok")
	}
}

func TestChain(t *testing.T) {
	t.Parallel()

	recordMiddleware := func(trace *[]string, name string) Middleware {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				*trace = append(*trace, name+":before")
				next.ServeHTTP(w, r)
				*trace = append(*trace, name+":after")
			})
		}
	}

	shortCircuitMiddleware := func(trace *[]string, name string, status int) Middleware {
		return func(_ http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				*trace = append(*trace, name+":abort")
				w.WriteHeader(status)
			})
		}
	}

	tests := []struct {
		name       string
		build      func(trace *[]string) []Middleware
		wantStatus int
		wantTrace  []string
	}{
		{
			name: "no middlewares",
			build: func(_ *[]string) []Middleware {
				return nil
			},
			wantStatus: http.StatusOK,
			wantTrace:  []string{"handler"},
		},
		{
			name: "pass through only",
			build: func(_ *[]string) []Middleware {
				return []Middleware{PassThrough}
			},
			wantStatus: http.StatusOK,
			wantTrace:  []string{"handler"},
		},
		{
			name: "preserves outer to inner execution order",
			build: func(trace *[]string) []Middleware {
				return []Middleware{
					recordMiddleware(trace, "first"),
					PassThrough,
					recordMiddleware(trace, "second"),
					recordMiddleware(trace, "third"),
				}
			},
			wantStatus: http.StatusOK,
			wantTrace: []string{
				"first:before",
				"second:before",
				"third:before",
				"handler",
				"third:after",
				"second:after",
				"first:after",
			},
		},
		{
			name: "ignores nil middleware entries",
			build: func(trace *[]string) []Middleware {
				return []Middleware{
					nil,
					recordMiddleware(trace, "only"),
					nil,
				}
			},
			wantStatus: http.StatusOK,
			wantTrace:  []string{"only:before", "handler", "only:after"},
		},
		{
			name: "supports short circuiting",
			build: func(trace *[]string) []Middleware {
				return []Middleware{
					recordMiddleware(trace, "first"),
					shortCircuitMiddleware(trace, "auth", http.StatusUnauthorized),
					recordMiddleware(trace, "never"),
				}
			},
			wantStatus: http.StatusUnauthorized,
			wantTrace:  []string{"first:before", "auth:abort", "first:after"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var trace []string
			handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				trace = append(trace, "handler")
				w.WriteHeader(http.StatusOK)
			})

			chained := Chain(tc.build(&trace)...)(handler)
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/healthz", nil)

			chained.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatus)
			}
			if !reflect.DeepEqual(trace, tc.wantTrace) {
				t.Errorf("trace = %v, want %v", trace, tc.wantTrace)
			}
		})
	}
}
