package rapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Request gathers all information about request
type Request struct {
	ID          int64  // storing record ID from URL
	Root        string // default JSON root key
	action      string
	req         *http.Request
	params      map[string]interface{}
	renderJSON  func(int, JSONData)
	render      func(int, string)
	renderError func(int, string)
	loadFile    func(string, string) (string, error)
}

// Init initializing controller
func (r *Request) Init(w http.ResponseWriter, req *http.Request, root, prefix string) {
	r.req = req
	r.Root = root

	urlParts := urlParts(strings.TrimPrefix(req.URL.Path, prefix))
	r.ID, _ = strconv.ParseInt(urlParts[0], 10, 64)
	r.action = r.makeAction(req.Method, urlParts)

	if strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
		r.renderJSON = func(code int, s JSONData) {
			RenderJSONgzip(w, code, s)
		}
	} else {
		r.renderJSON = func(code int, s JSONData) {
			RenderJSON(w, code, s)
		}
	}

	r.render = func(code int, s string) {
		w.WriteHeader(code)
		w.Write([]byte(s))
	}

	r.renderError = func(code int, s string) {
		http.Error(w, s, code)
	}

	r.loadFile = func(field, dir string) (string, error) {
		r.req.ParseMultipartForm(32 << 20)
		file, handler, err := r.req.FormFile(field)
		if err != nil {
			return "", err
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile(dir+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return "", err
		}
		defer f.Close()
		io.Copy(f, file)
		return handler.Filename, nil
	}

	r.params = make(map[string]interface{})
}

func (r *Request) makeAction(method string, urlParts map[int]string) string {
	if len(urlParts[1]) > 0 {
		return method + capitalize(urlParts[1])
	}
	if r.ID > 0 {
		switch method {
		case "GET":
			return "Show"
		case "POST", "PUT":
			return "Update"
		case "DELETE":
			return "Destroy"
		}
	}

	if len(urlParts[0]) > 0 {
		return method + capitalize(urlParts[0])
	}

	switch method {
	case "GET":
		return "Index"
	case "POST":
		return "Create"
	}

	return "WrongAction"
}

// LoadJSONRequest extracting JSON request by key
// from request body into interface
func (r *Request) LoadJSONRequest(root string, v interface{}) {
	if root == "" {
		extractJSONPayload(r.req.Body, &v)
		return
	}

	var s []byte
	var body JSONData
	extractJSONPayload(r.req.Body, &body)
	s, _ = json.Marshal(body[root])
	json.Unmarshal(s, &v)
}

// QueryParam returns URL query param
func (r *Request) QueryParam(s string) string {
	return r.req.URL.Query().Get(s)
}

// SetParam set custom parameter for current request
func (r *Request) SetParam(k string, v interface{}) {
	r.params[k] = v
}

// Param returns custom parameter for current request
func (r *Request) Param(k string) interface{} {
	return r.params[k]
}

// Header returns request header
func (r *Request) Header(s string) string {
	return r.req.Header.Get(s)
}

// CurrentAction returns current controller action
func (r *Request) CurrentAction() string {
	return r.action
}

// RenderJSON rendering JSON to client
func (r *Request) RenderJSON(code int, s JSONData) {
	r.renderJSON(code, s)
}

// RenderJSONError rendering error to client in JSON format
func (r *Request) RenderJSONError(code int, s string) {
	r.renderJSON(code, JSONData{"errors": JSONData{"message": []string{s}}})
}

// Render rendering string to client
func (r *Request) RenderString(code int, s string) {
	r.render(code, s)
}

// RenderError rendering error to client
func (r *Request) RenderError(code int, s string) {
	r.renderError(code, s)
}

// LoadFile handling file uploads
func (r *Request) LoadFile(field, dir string) (string, error) {
	return r.loadFile(field, dir)
}
