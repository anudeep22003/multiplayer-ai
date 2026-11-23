package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/anudeep22003/conman/internal"
)

const shutdownTimeout = 15 * time.Second

// New configures the HTTP listener and timeouts for handler.
func New(cfg *internal.Config, handler http.Handler) *http.Server {
	return &http.Server{
		// Addr optionally specifies the TCP address for the server to listen on,
		// in the form "host:port". If empty, ":http" (port 80) is used.
		// The service names are defined in RFC 6335 and assigned by IANA.
		// See net.Dial for details of the address format.
		Addr:    net.JoinHostPort(cfg.Host, strconv.Itoa(int(cfg.Port))),
		Handler: handler,

		// ReadTimeout is the maximum duration for reading the entire
		// request, including the body. A zero or negative value means
		// there will be no timeout.
		//
		// Because ReadTimeout does not let Handlers make per-request
		// decisions on each request body's acceptable deadline or
		// upload rate, most users will prefer to use
		// ReadHeaderTimeout. It is valid to use them both.
		ReadTimeout: cfg.ReadTimeout,

		// WriteTimeout is the maximum duration before timing out
		// writes of the response. It is reset whenever a new
		// request's header is read. Like ReadTimeout, it does not
		// let Handlers make decisions on a per-request basis.
		// A zero or negative value means there will be no timeout.
		WriteTimeout: cfg.WriteTimeout,

		// IdleTimeout is the maximum amount of time to wait for the
		// next request when keep-alives are enabled. If zero, the value
		// of ReadTimeout is used. If negative, or if zero and ReadTimeout
		// is zero or negative, there is no timeout.
		IdleTimeout: cfg.IdleTimeout,
	}
}

// Run serves HTTP and waits for a graceful shutdown when ctx is canceled.
func Run(ctx context.Context, srv *http.Server) error {
	serveErrors := make(chan error, 1)
	go func() {
		serveErrors <- srv.ListenAndServe()
	}()

	var serveErr error
	select {
	case serveErr = <-serveErrors:
	case <-ctx.Done():
		// The signal context is already canceled. Give active requests their own deadline.
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}
		serveErr = <-serveErrors
	}

	if errors.Is(serveErr, http.ErrServerClosed) {
		return nil
	}
	return fmt.Errorf("serve HTTP: %w", serveErr)
}
