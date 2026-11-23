package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/anudeep22003/conman/internal"
)

func TestNew(t *testing.T) {
	t.Parallel()

	cfg := internal.Config{
		Host:         "127.0.0.1",
		Port:         9090,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	handler := http.NewServeMux()
	srv := New(&cfg, handler)
	if srv.Handler != handler {
		t.Fatal("server did not use supplied handler")
	}

	if srv.Addr != "127.0.0.1:9090" {
		t.Errorf("Addr = %q, want %q", srv.Addr, "127.0.0.1:9090")
	}
	if srv.ReadTimeout != cfg.ReadTimeout || srv.WriteTimeout != cfg.WriteTimeout || srv.IdleTimeout != cfg.IdleTimeout {
		t.Errorf("timeouts = (%v, %v, %v), want (%v, %v, %v)",
			srv.ReadTimeout, srv.WriteTimeout, srv.IdleTimeout,
			cfg.ReadTimeout, cfg.WriteTimeout, cfg.IdleTimeout,
		)
	}
}
