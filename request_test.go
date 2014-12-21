package rapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"
)

var httpWriter http.ResponseWriter

// newRequest is a helper function to create a new request with a method and url
func newRequest(method, url string, body string) *http.Request {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-API-Token", "token1")
	return req
}

func newRecorder() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

func assertEqual(t *testing.T, expect interface{}, v interface{}) {
	if v != expect {
		_, fname, lineno, ok := runtime.Caller(1)
		if !ok {
			fname, lineno = "<UNKNOWN>", -1
		}
		t.Errorf("FAIL: %s:%d\nExpected: %#v\nReceived: %#v", fname, lineno, expect, v)
	}
}

func assertNotEqual(t *testing.T, expect interface{}, v interface{}) {
	if v == expect {
		_, fname, lineno, ok := runtime.Caller(1)
		if !ok {
			fname, lineno = "<UNKNOWN>", -1
		}
		t.Errorf("FAIL: %s:%d\nExpected: %#v\nReceived: %#v", fname, lineno, expect, v)
	}
}

func newReq(w http.ResponseWriter, req *http.Request, root, prefix string) *Request {
	r := &Request{}
	r.Init(w, req, root, prefix)
	return r
}

func TestMakeAction(t *testing.T) {
	p := "/pages"
	r := newReq(httpWriter, newRequest("GET", "http://localhost/pages/10", "{}"), "root", p)
	assertEqual(t, "Show", r.Action)
	assertEqual(t, "10", r.URL.ID)
	assertEqual(t, int64(10), r.URL.ID64())
	assertEqual(t, "", r.URL.Action)

	r = newReq(httpWriter, newRequest("GET", "http://localhost/pages/10/edit", "{}"), "root", p)
	assertEqual(t, "GETEdit", r.Action)
	assertEqual(t, "10", r.URL.ID)
	assertEqual(t, int64(10), r.URL.ID64())
	assertEqual(t, "edit", r.URL.Action)

	r = newReq(httpWriter, newRequest("POST", "http://localhost/pages/10", "{}"), "root", p)
	assertEqual(t, "Update", r.Action)
	assertEqual(t, "10", r.URL.ID)
	assertEqual(t, int64(10), r.URL.ID64())
	assertEqual(t, "", r.URL.Action)

	r = newReq(httpWriter, newRequest("POST", "http://localhost/pages/10/edit", "{}"), "root", p)
	assertEqual(t, "POSTEdit", r.Action)
	assertEqual(t, "10", r.URL.ID)
	assertEqual(t, int64(10), r.URL.ID64())
	assertEqual(t, "edit", r.URL.Action)

	r = newReq(httpWriter, newRequest("PUT", "http://localhost/pages/10", "{}"), "root", p)
	assertEqual(t, "Update", r.Action)
	assertEqual(t, "10", r.URL.ID)
	assertEqual(t, int64(10), r.URL.ID64())
	assertEqual(t, "", r.URL.Action)

	r = newReq(httpWriter, newRequest("PUT", "http://localhost/pages/10/edit", "{}"), "root", p)
	assertEqual(t, "PUTEdit", r.Action)
	assertEqual(t, "10", r.URL.ID)
	assertEqual(t, int64(10), r.URL.ID64())
	assertEqual(t, "edit", r.URL.Action)

	r = newReq(httpWriter, newRequest("DELETE", "http://localhost/pages/10", "{}"), "root", p)
	assertEqual(t, "Destroy", r.Action)
	assertEqual(t, "10", r.URL.ID)
	assertEqual(t, int64(10), r.URL.ID64())
	assertEqual(t, "", r.URL.Action)

	r = newReq(httpWriter, newRequest("DELETE", "http://localhost/pages/10/edit", "{}"), "root", p)
	assertEqual(t, "DELETEEdit", r.Action)
	assertEqual(t, "10", r.URL.ID)
	assertEqual(t, int64(10), r.URL.ID64())
	assertEqual(t, "edit", r.URL.Action)
}

func TestQueryParams(t *testing.T) {
	req := newRequest("GET", "http://localhost/?p1=1&p2=2", "{}")
	r := Request{}
	r.Init(httpWriter, req, "root", "")
	assertEqual(t, "1", r.QueryParam("p1"))
	assertEqual(t, "2", r.QueryParam("p2"))
	assertEqual(t, "", r.QueryParam("p3"))
}

func TestHeader(t *testing.T) {
	req := newRequest("GET", "http://localhost", "{}")
	r := Request{}
	r.Init(httpWriter, req, "root", "")
	assertEqual(t, "token1", r.Header("X-API-Token"))
	assertEqual(t, "", r.Header("X-API-Token1"))
}

func TestBody(t *testing.T) {
	req := newRequest("GET", "http://localhost/", "{\"id\":2}")
	r := Request{}
	r.Init(httpWriter, req, "root", "")
	var res interface{}
	res = nil
	r.LoadJSONRequest("", &res)
	in := fmt.Sprintf("%#v", res)
	out := fmt.Sprintf("%#v", map[string]interface{}{"id": 2})
	assertEqual(t, out, in)

	req = newRequest("GET", "http://localhost/", "{\"id\":2}")
	r = Request{}
	r.Init(httpWriter, req, "root", "")
	res = nil
	r.LoadJSONRequest("id", &res)
	in = fmt.Sprintf("%#v", res)
	assertEqual(t, "2", in)

	req = newRequest("GET", "http://localhost/", "{\"id\":2}")
	r = Request{}
	r.Init(httpWriter, req, "root", "")
	res = nil
	r.LoadJSONRequest("id1", &res)
	assertEqual(t, nil, res)
}

type TestC struct {
	Request
}

func (t *TestC) Index() {
	t.RenderJSON(200, JSONData{t.Root: "index"})
}

func (t *TestC) Show() {
	t.RenderJSON(200, JSONData{t.Root: "show"})
}

func (t *TestC) Create() {
	var i interface{}
	t.LoadJSONRequest("root", &i)
	t.RenderJSON(200, JSONData{t.Root: i})
}

func TestReponseIndex(t *testing.T) {
	req := newRequest("GET", "http://localhost/pages/", "{}")
	handler := handle(&TestC{}, "page", "/pages")
	rec := newRecorder()
	handler(rec, req)
	assertEqual(t, "{\"page\":\"index\"}\n", string(rec.Body.Bytes()))
}

func TestReponseShow(t *testing.T) {
	req := newRequest("GET", "http://localhost/pages/10", "{}")
	handler := handle(&TestC{}, "page", "/pages")
	rec := newRecorder()
	handler(rec, req)
	assertEqual(t, "{\"page\":\"show\"}\n", string(rec.Body.Bytes()))
}

func TestReponseCreate(t *testing.T) {
	req := newRequest("POST", "http://localhost/pages", `{"root":[{"id":1}]}`)
	handler := handle(&TestC{}, "page", "/pages")
	rec := newRecorder()
	handler(rec, req)
	assertEqual(t, "{\"page\":[{\"id\":1}]}\n", string(rec.Body.Bytes()))
}

func BenchmarkHandleIndex(b *testing.B) {
	req := newRequest("GET", "http://localhost/pages/", "{}")
	handler := handle(&TestC{}, "page", "/pages")

	for n := 0; n < b.N; n++ {
		handler(newRecorder(), req)
	}
}

func BenchmarkHandleShow(b *testing.B) {
	req := newRequest("GET", "http://localhost/pages/10", "{}")
	handler := handle(&TestC{}, "page", "/pages")

	for n := 0; n < b.N; n++ {
		handler(newRecorder(), req)
	}
}

func BenchmarkHandleCreate(b *testing.B) {
	req := newRequest("POST", "http://localhost/pages/", `{"root":[{"id":1}]}`)
	handler := handle(&TestC{}, "page", "/pages")

	for n := 0; n < b.N; n++ {
		handler(newRecorder(), req)
	}
}
