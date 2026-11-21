package server

import "net/http"

// Middleware wraps an http.Handler with cross-cutting request/response logic.
type Middleware func(http.Handler) http.Handler

// Chain composes middlewares so they execute in the order passed:
// Chain(m1, m2)(h) runs m1 first, then m2, and finally h.
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			if middlewares[i] != nil {
				next = middlewares[i](next)
			}
		}
		return next
	}
}

// PassThrough forwards the request to the next handler unchanged.
func PassThrough(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
