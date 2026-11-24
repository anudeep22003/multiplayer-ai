package handlers

import (
	"net/http"
	"strings"
)

type handler interface {
	RegisterRoutes() *http.ServeMux
}

type RouteGroup struct {
	Prefix  string
	Handler http.Handler
}

func Assemble(root *http.ServeMux, rgs ...RouteGroup) *http.ServeMux {
	// given a root and routegroup, assemble all the muxes onto it

	for _, r := range rgs {
		mount(root, r.Prefix, r.Handler)
	}

	return root

}

func mount(root *http.ServeMux, prefix string, h http.Handler) {
	// safe handling somebody adding a `/` at the end
	prefix = strings.TrimSuffix(prefix, "/")
	root.Handle(prefix+"/", http.StripPrefix(prefix, h))
}
