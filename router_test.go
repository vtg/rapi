package rapi

import (
	"fmt"
	"net/http"
	"testing"
)

func HandlerForTest(w http.ResponseWriter, req *http.Request)  {}
func HandlerForTest1(w http.ResponseWriter, req *http.Request) {}

func TestHandleFunc(t *testing.T) {
	r := NewRouter()

	r.HandleFunc("/pages", HandlerForTest)
	assertEqual(t, 0, len(r.routes))
	assertEqual(t, 1, len(r.namedRoutes))

	assertEqual(t, fmt.Sprint(http.NotFoundHandler()), fmt.Sprint(r.match("/pages/1")))
	assertEqual(t, fmt.Sprint(HandlerForTest), fmt.Sprint(r.match("/pages")))
}

func TestPathPrefix(t *testing.T) {
	r := NewRouter()
	p := r.PathPrefix("/api")
	p.HandleFunc("//pages2", HandlerForTest)
	p.HandleFunc("/pages1", HandlerForTest1)
	assertEqual(t, 2, len(r.routes))
	assertEqual(t, 0, len(r.namedRoutes))

	assertEqual(t, fmt.Sprint(http.NotFoundHandler()), fmt.Sprint(r.match("/pages1/1")))
	assertEqual(t, fmt.Sprint(HandlerForTest1), fmt.Sprint(r.match("/api/pages1/1")))

	assertEqual(t, fmt.Sprint(http.NotFoundHandler()), fmt.Sprint(r.match("/pages2/1")))
	assertEqual(t, fmt.Sprint(HandlerForTest), fmt.Sprint(r.match("/api/pages2/1")))
}

func TestRoute(t *testing.T) {

	type C struct{ Request }

	r := NewRouter()

	c := &C{}

	r.Route("/pages", c, "c")
	assertEqual(t, 1, len(r.routes))
	assertEqual(t, 0, len(r.namedRoutes))

	assertEqual(t, "", r.matchNamed("/pages/1"))
	assertEqual(t, "/pages", r.matchCommon("/pages/1"))

	p := r.PathPrefix("/api")
	p.Route("/pages", c, "c")
	assertEqual(t, 2, len(r.routes))
	assertEqual(t, 0, len(r.namedRoutes))

	assertEqual(t, "", r.matchNamed("/api/pages/2"))
	assertEqual(t, "/api/pages", r.matchCommon("/api/pages/2"))

}

func TestRoutesOrder(t *testing.T) {

	type C struct {
		Request
	}

	r := NewRouter()
	r.HandleFunc("/a", HandlerForTest)
	r.Route("/aa", &C{}, "c")
	r.Route("/aaa", &C{}, "c")
	r.Route("/aaaa", &C{}, "c")
	r.Route("/aaaaa", &C{}, "c")
	r.Route("/a/a", &C{}, "c")

	assertEqual(t, "/aa", r.matchCommon("/aa/1"))
	assertEqual(t, "/a", r.matchNamed("/a"))
	assertEqual(t, "/aaa", r.matchCommon("/aaa/"))
	assertEqual(t, "/aaaa", r.matchCommon("/aaaa/22"))
	assertEqual(t, "/a/a", r.matchCommon("/a/a/"))
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
