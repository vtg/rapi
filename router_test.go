package rapi

import (
	"net/http"
	"testing"
)

func HandlerForTest(w http.ResponseWriter, req *http.Request) {

}

func TestHandleFunc(t *testing.T) {
	r := NewRouter()

	r.HandleFunc("/pages", HandlerForTest)
	assertEqual(t, 0, len(r.routes))
	assertEqual(t, 1, len(r.namedRoutes))

	route := r.match("/pages/1")
	assertEqual(t, nil, route.handler)

	route = r.match("/pages")
	assertEqual(t, "/pages", route.prefix)
	assertNotEqual(t, nil, route.handler)
}

func TestPathPrefix(t *testing.T) {
	r := NewRouter()
	p := r.PathPrefix("/api")
	p.HandleFunc("//pages2", HandlerForTest)
	p.HandleFunc("/pages1", HandlerForTest)
	assertEqual(t, 2, len(r.routes))
	assertEqual(t, 0, len(r.namedRoutes))

	route := r.match("/pages1/1")
	assertEqual(t, nil, route.handler)

	route = r.match("/api/pages1/1")
	assertEqual(t, "/api/pages1", route.prefix)
	assertEqual(t, "/1", route.match)

	route = r.match("/pages2/1")
	assertEqual(t, nil, route.handler)

	route = r.match("/api/pages2/1")
	assertEqual(t, "/api/pages2", route.prefix)
	assertEqual(t, "/1", route.match)
}

func TestRoute(t *testing.T) {

	type C struct {
		Request
	}

	r := NewRouter()

	r.Route("/pages", &C{}, "c")
	assertEqual(t, 1, len(r.routes))
	assertEqual(t, 0, len(r.namedRoutes))

	route := r.match("/pages/1")
	assertEqual(t, "/pages", route.prefix)
	assertEqual(t, "/1", route.match)

	p := r.PathPrefix("/api")
	p.Route("/pages", &C{}, "c")
	assertEqual(t, 2, len(r.routes))
	assertEqual(t, 0, len(r.namedRoutes))

	route = r.match("/api/pages/2")
	assertEqual(t, "/api/pages", route.prefix)
	assertEqual(t, "/2", route.match)

}

func BenchmarkMatch(b *testing.B) {
	r := NewRouter()
	p := r.PathPrefix("/api")
	p.HandleFunc("/pages1", HandlerForTest)
	p.HandleFunc("/pages2", HandlerForTest)
	p.HandleFunc("/pages3", HandlerForTest)
	p.HandleFunc("/pages4", HandlerForTest)
	p.HandleFunc("/pages5", HandlerForTest)
	p.HandleFunc("/pages6", HandlerForTest)

	for n := 0; n < b.N; n++ {
		r.match("/api/pages1/1")
	}
}
