package rapi

import (
	"net/http"
	"strings"
)

type Router struct {
	namedRoutes map[string]*Route
	routes      []*Route
}

func NewRouter() *Router {
	return &Router{namedRoutes: make(map[string]*Route)}
}

// HandleFunc registers a new route with a matcher for the URL path.
// See Route.HandlerFunc().
func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) {
	r.NewRoute("").HandleFunc(path, f)
}

// Route registers a new route with a matcher for URL path
// and registering controller handler
func (r *Router) Route(path string, i Controller, rootKey string, funcs ...ReqFunc) {
	route := r.NewRoute(path)
	route.HandlerFunc(handle(i, rootKey, route.prefix, funcs...)).addRoute(false)
}

// HandlePrefix registers a new handler to serve prefix
func (r *Router) HandlePrefix(path string, handler http.Handler) {
	r.NewRoute(path).Handler(handler).addRoute(false)
}

// NewRoute registers an empty route.
func (r *Router) NewRoute(prefix string) *Route {
	return &Route{router: r, prefix: prefix}
}

func (r *Router) addRoute(rt *Route) {
	r.routes = append(r.routes, rt)
}

func (r *Router) addNamedRoute(rt *Route) {
	r.namedRoutes[rt.prefix] = rt
}

func (r *Router) match(path string) *Route {
	for k := range r.namedRoutes {
		if k == path {
			return r.namedRoutes[k]
		}
	}
	for k := range r.routes {
		prefix := r.routes[k].prefix
		if strings.HasPrefix(path, prefix) {
			return r.routes[k]
		}
	}
	return &Route{}
}

func (r *Router) PathPrefix(s string) *Route {
	return r.NewRoute(s)
}

// ServeHTTP dispatches the handler registered in the matched route.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if p := cleanPath(req.URL.Path); p != req.URL.Path {
		url := *req.URL
		url.Path = p
		p = url.String()

		w.Header().Set("Location", p)
		w.WriteHeader(http.StatusMovedPermanently)
		return
	}

	var handler http.Handler
	route := r.match(req.URL.Path)
	if route.handler != nil {
		handler = route.handler
	} else {
		handler = http.NotFoundHandler()
	}

	handler.ServeHTTP(w, req)
}
