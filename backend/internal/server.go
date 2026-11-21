package internal

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
)

// NewServer configures an HTTP server with the backend's routes and timeouts.
func NewServer(cfg *Config) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{"Hello": "World"}

		bytes, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "json parsing error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write(bytes)

	})

	return &http.Server{
		// Addr optionally specifies the TCP address for the server to listen on,
		// in the form "host:port". If empty, ":http" (port 80) is used.
		// The service names are defined in RFC 6335 and assigned by IANA.
		// See net.Dial for details of the address format.
		Addr:    net.JoinHostPort(cfg.Host, strconv.Itoa(int(cfg.Port))),
		Handler: mux, // handler to invoke, http.DefaultServeMux if nil

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
