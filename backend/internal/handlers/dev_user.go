package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

type devUserIDKey struct{}

// DevUserMiddleware attaches the optional development user ID to the request.
func DevUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.Header.Values("X-Dev-User-Id")
		if len(values) > 0 {
			var id pgtype.UUID
			if len(values) != 1 || id.Scan(values[0]) != nil || !strings.EqualFold(values[0], id.String()) {
				http.Error(w, "invalid X-Dev-User-Id", http.StatusBadRequest)
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), devUserIDKey{}, id.String()))
		}
		next.ServeHTTP(w, r)
	})
}
