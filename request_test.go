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

func TestMakeAction(t *testing.T) {
	r := Request{}

	r.ID = 10
	act := r.makeAction("GET", urlParts("/10"))
	assertEqual(t, "Show", act)
	act = r.makeAction("POST", urlParts("/10"))
	assertEqual(t, "Update", act)
	act = r.makeAction("PUT", urlParts("/10"))
	assertEqual(t, "Update", act)
	act = r.makeAction("DELETE", urlParts("/10"))
	assertEqual(t, "Destroy", act)
	act = r.makeAction("GET", urlParts("/10/edit"))
	assertEqual(t, "GETEdit", act)
	act = r.makeAction("POST", urlParts("/10/edit"))
	assertEqual(t, "POSTEdit", act)
	act = r.makeAction("PUT", urlParts("/10/edit"))
	assertEqual(t, "PUTEdit", act)
	act = r.makeAction("DELETE", urlParts("/10/edit"))
	assertEqual(t, "DELETEEdit", act)

	r.ID = 0
	act = r.makeAction("GET", urlParts("/action"))
	assertEqual(t, "GETAction", act)
	act = r.makeAction("POST", urlParts("/action"))
	assertEqual(t, "POSTAction", act)
	act = r.makeAction("PUT", urlParts("/action"))
	assertEqual(t, "PUTAction", act)
	act = r.makeAction("DELETE", urlParts("/action"))
	assertEqual(t, "DELETEAction", act)
}

func TestQueryParams(t *testing.T) {
	req := newRequest("GET", "http://localhost/?p1=1&p2=2", "{}")
	r := Request{}
	r.Init(httpWriter, req, "root")
	assertEqual(t, "1", r.QueryParam("p1"))
	assertEqual(t, "2", r.QueryParam("p2"))
	assertEqual(t, "", r.QueryParam("p3"))
}

func TestHeader(t *testing.T) {
	req := newRequest("GET", "http://localhost", "{}")
	r := Request{}
	r.Init(httpWriter, req, "root")
	assertEqual(t, "token1", r.Header("X-API-Token"))
	assertEqual(t, "", r.Header("X-API-Token1"))
}

func TestBody(t *testing.T) {
	req := newRequest("GET", "http://localhost/", "{\"id\":2}")
	r := Request{}
	r.Init(httpWriter, req, "root")
	var res interface{}
	res = nil
	r.LoadJSONRequest("", &res)
	in := fmt.Sprintf("%#v", res)
	out := fmt.Sprintf("%#v", map[string]interface{}{"id": 2})
	assertEqual(t, out, in)

	req = newRequest("GET", "http://localhost/", "{\"id\":2}")
	r = Request{}
	r.Init(httpWriter, req, "root")
	res = nil
	r.LoadJSONRequest("id", &res)
	in = fmt.Sprintf("%#v", res)
	assertEqual(t, "2", in)

	req = newRequest("GET", "http://localhost/", "{\"id\":2}")
	r = Request{}
	r.Init(httpWriter, req, "root")
	res = nil
	r.LoadJSONRequest("id1", &res)
	assertEqual(t, nil, res)
}
