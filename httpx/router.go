package httpx

import (
	"hybrid-srv/core"
)

// TODO 改trie树
type Router struct {
	routes map[string]map[string]core.HandlerFunc
}

func NewRouter() *Router {
	return &Router{routes: make(map[string]map[string]core.HandlerFunc)}
}

func (r *Router) Add(method, path string, handler core.HandlerFunc) {
	if r.routes[method] == nil {
		r.routes[method] = make(map[string]core.HandlerFunc)
	}
	r.routes[method][path] = handler
}

func (r *Router) Match(method, path string) core.HandlerFunc {
	if r.routes[method] == nil {
		return nil
	}
	return r.routes[method][path]
}
