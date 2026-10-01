package router

import (
	"net/http"

	"github.com/stevanusy21/golang_sandbox/pkg/middleware"
)

type Router struct {
	Mux *http.ServeMux
}

func NewRouter(mux *http.ServeMux) *Router {
	return &Router{Mux: mux}
}

func (r *Router) Public(pattern string, handler http.HandlerFunc) {
	r.Mux.HandleFunc(pattern, handler)
}

func (r *Router) Protected(pattern string, handler http.HandlerFunc) {
	r.Mux.Handle(pattern, middleware.AuthMiddleware(handler))
}
