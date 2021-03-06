package rapi

import "net/http"

type Route struct {
	router  *Router
	prefix  string
	handler http.Handler
}

// HandleFunc setting function to handle route
func (r *Route) HandleFunc(s string, f func(http.ResponseWriter, *http.Request)) {
	r.NewRoute(s).HandlerFunc(f).addRoute(r.prefix == "")
}

// Route registers a new route with a matcher for URL path
// and registering controller handler
// ex:
//    r := api.NewRouter()
//    api = r.PathPrefix("/api/v1")
//    api.Route("/pages", &PagesController{}, "page", AuthFunc)
// where
//  - PagesController is the type implementing Controller
//  - "page" is the root key for json request/response
//  - AuthFunc is middleware function that implements ReqFunc.
//
func (r *Route) Route(path string, i Controller, rootKey string, funcs ...ReqFunc) {
	rt := r.NewRoute(path)
	rt.HandlerFunc(handle(i, rootKey, rt.prefix, implements(i), funcs...)).addRoute(false)
}

// FileServer provides static files serving
// ex:
//    r := api.NewRouter()
//    dirIndex := false
//    preferGzip := false
//    r.PathPrefix("/images/").FileServer("./public", dirIndex, preferGzip)
//
// where
//  - dirIndex specifying if it should display directory content or not
//  - preferGzip specifying if it should look for gzipped file version
//
func (r *Route) FileServer(path string, b ...bool) {
	r.Handler(fileServer(path, b)).addRoute(false)
}

// NewRoute registers an empty route.
func (r *Route) NewRoute(prefix string) *Route {
	return &Route{router: r.router, prefix: cleanPath(r.prefix + prefix)}
}

// Handler sets a handler for the route.
func (r *Route) Handler(handler http.Handler) *Route {
	r.handler = handler
	return r
}

// HandlerFunc sets a handler function for the route.
func (r *Route) HandlerFunc(f func(http.ResponseWriter, *http.Request)) *Route {
	return r.Handler(http.HandlerFunc(f))
}

func (r *Route) addRoute(named bool) *Route {
	if named {
		r.router.addNamedRoute(r)
	} else {
		r.router.addRoute(r)
	}
	return r
}
