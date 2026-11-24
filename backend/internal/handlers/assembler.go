package handlers

import (
	"net/http"
	"strings"

	"github.com/anudeep22003/conman/internal/handlers/health"
)

type handler interface {
	RegisterRoutes() *http.ServeMux
}

type RouteGroup struct {
	Prefix  string
	Handler http.Handler
}

type deps struct {
	dbPool DB
}

func AssembledHttpMux(deps *deps) *http.ServeMux {
	root := http.NewServeMux()

	routes := []RouteGroup{
		{
			Prefix:  "",
			Handler: health.RegisterRoutes(deps.dbPool),
		},
	}

	for _, r := range routes {
		mount(root, r.Prefix, r.Handler)
	}

	return root

}

func mount(root *http.ServeMux, prefix string, h http.Handler) {
	prefix = strings.TrimSuffix(prefix, "/")
	root.Handle(prefix+"/", http.StripPrefix(prefix, h))
}
