package rapi

import (
	"net/http"
	"strings"

	"github.com/gorilla/context"
)

type Router struct {
	namedRoutes map[string]*Route
	routes      []*Route
	KeepContext bool
}

func NewRouter() *Router {
	return &Router{namedRoutes: make(map[string]*Route), KeepContext: false}
}

// HandleFunc registers a new route with a matcher for the URL path.
// See Route.HandlerFunc().
func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *Route {
	return r.NewRoute("").HandleFunc(path, f)
}

// Route registers a new route with a matcher for URL path
// and registering controller handler
func (r *Router) Route(path string, i Controller, rootKey string, funcs ...ReqFunc) *Route {
	route := r.NewRoute(path).HandlerFunc(handle(i, rootKey, funcs))
	route.addRoute(false)
	return route
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
			r.routes[k].match = strings.TrimPrefix(path, prefix)
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
		setCurrentPath(req, route.match)
	} else {
		handler = http.NotFoundHandler()
	}

	if !r.KeepContext {
		defer context.Clear(req)
	}
	handler.ServeHTTP(w, req)
}

type contextKey int

const pKey contextKey = iota

func setCurrentPath(r *http.Request, s string) {
	context.Set(r, pKey, s)
}

// currentPath returns request URL without prefix
func currentPath(r *http.Request) string {
	if rv := context.Get(r, pKey); rv != nil {
		return rv.(string)
	}
	return ""
}
