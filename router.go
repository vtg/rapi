package rapi

import (
	"net/http"
	"reflect"
	"sort"
	"strings"
)

type Router struct {
	routes      map[string]http.Handler
	namedRoutes map[string]http.Handler
	keys        []string
	namedKeys   []string
}

func NewRouter() *Router {
	return &Router{namedRoutes: make(map[string]http.Handler), routes: make(map[string]http.Handler)}
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
	route.HandlerFunc(handle(i, rootKey, route.prefix, implements(i), funcs...)).addRoute(false)
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
	r.routes[rt.prefix] = rt.handler
	r.setKeys()
}

func (r *Router) addNamedRoute(rt *Route) {
	r.namedRoutes[rt.prefix] = rt.handler
	r.setNamedKeys()
}

func (r *Router) setKeys() {
	r.keys = make([]string, len(r.routes))
	for key := range r.routes {
		r.keys = append(r.keys, key)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(r.keys)))
}

func (r *Router) setNamedKeys() {
	r.namedKeys = make([]string, len(r.namedRoutes))
	for key := range r.namedRoutes {
		r.namedKeys = append(r.namedKeys, key)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(r.namedKeys)))
}

func (r *Router) matchNamed(path string) string {
	for _, v := range r.namedKeys {
		if v == path {
			return v
		}
	}
	return ""
}

func (r *Router) matchCommon(path string) string {
	for _, v := range r.keys {
		if strings.HasPrefix(path, v) {
			return v
		}
	}
	return ""
}

func (r *Router) match(path string) http.Handler {
	if k := r.matchNamed(path); k != "" {
		return r.namedRoutes[k]
	}

	if k := r.matchCommon(path); k != "" {
		return r.routes[k]
	}

	return http.NotFoundHandler()
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

	r.match(req.URL.Path).ServeHTTP(w, req)
}

var meths = []string{"GET", "POST"}

// implements extracting custom methods from controller
// custom method names should begin from GET or POST
func implements(v interface{}) []string {
	res := []string{}
	t := reflect.TypeOf(v)
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		for _, v := range meths {
			if strings.HasPrefix(m.Name, v) {
				res = append(res, m.Name)
				continue
			}
		}
	}
	return res
}
