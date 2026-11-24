package handlers

import (
	"net/http"
	"strings"
)

type RouteGroup struct {
	Prefix  string
	Handler http.Handler
}

func Assemble(rgs ...RouteGroup) *http.ServeMux {
	root := http.NewServeMux()
	for _, r := range rgs {
		mount(root, r.Prefix, r.Handler)
	}
	return root
}

func mount(root *http.ServeMux, prefix string, h http.Handler) {
	prefix = strings.TrimSuffix(prefix, "/")
	if prefix == "" {
		root.Handle("/", h)
		return
	}
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	root.Handle(prefix+"/", http.StripPrefix(prefix, h))
}
