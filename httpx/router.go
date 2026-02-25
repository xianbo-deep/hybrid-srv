package httpx

import (
	"hybrid-srv/core"
)

type Router struct {
	routes map[string]map[string]core.Handler
}

func NewRouter() *Router {
	return &Router{routes: make(map[string]map[string]core.Handler)}
}

func (r *Router) Add(method, path string, handler core.Handler) {
	if r.routes[method] == nil {
		r.routes[method] = make(map[string]core.Handler)
	}
	r.routes[method][path] = handler
}

func (r *Router) Match(method, path string) core.Handler {
	if r.routes[method] == nil {
		return nil
	}
	return r.routes[method][path]
}
