package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDevUserMiddleware(t *testing.T) {
	const id = "a035565a-896f-4fdf-9b5a-185e78fb0128"
	for _, tc := range []struct {
		headers []string
		wantID  any
		status  int
	}{
		{nil, nil, 200}, {[]string{id}, id, 200}, {[]string{""}, nil, 400}, {[]string{"bad"}, nil, 400}, {[]string{id, id}, nil, 400},
	} {
		req := httptest.NewRequest("GET", "/health", nil)
		for _, value := range tc.headers {
			req.Header.Add("X-Dev-User-Id", value)
		}
		rec := httptest.NewRecorder()
		DevUserMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got := r.Context().Value(devUserIDKey{}); got != tc.wantID {
				t.Errorf("context ID = %v, want %v", got, tc.wantID)
			}
		})).ServeHTTP(rec, req)
		if rec.Code != tc.status {
			t.Errorf("headers %q: status = %d, want %d", tc.headers, rec.Code, tc.status)
		}
	}
}
