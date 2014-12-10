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
	ID     int64  // storing record ID from URL
	Root   string // default JSON root key
	Action string
	params map[string]interface{}

	req *http.Request
	w   http.ResponseWriter
}

// Init initializing controller
func (r *Request) Init(w http.ResponseWriter, req *http.Request, root, prefix string) {
	r.w = w
	r.req = req
	r.Root = root

	urlParts := getParts(req.URL.Path, prefix)
	r.ID = urlParts.ID()
	r.Action = r.makeAction(urlParts)

	r.params = make(map[string]interface{})
}

func (r *Request) makeAction(urlParts parts) string {
	if urlParts.action != "" {
		return r.req.Method + capitalize(urlParts.action)
	}
	if r.ID > 0 {
		switch r.req.Method {
		case "GET":
			return "Show"
		case "POST", "PUT":
			return "Update"
		case "DELETE":
			return "Destroy"
		}
	}

	if urlParts.id != "" {
		return r.req.Method + capitalize(urlParts.id)
	}

	switch r.req.Method {
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
	defer r.req.Body.Close()

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
	return r.Action
}

// RenderJSON rendering JSON to client
func (r *Request) RenderJSON(code int, s JSONData) {
	if strings.Contains(r.req.Header.Get("Accept-Encoding"), "gzip") {
		RenderJSONgzip(r.w, code, s)
	} else {
		RenderJSON(r.w, code, s)
	}
}

// RenderJSONError rendering error to client in JSON format
func (r *Request) RenderJSONError(code int, s string) {
	r.RenderJSON(code, JSONData{"errors": JSONData{"message": []string{s}}})
}

// Render rendering string to client
func (r *Request) RenderString(code int, s string) {
	r.w.WriteHeader(code)
	r.w.Write([]byte(s))
}

// RenderError rendering error to client
func (r *Request) RenderError(code int, s string) {
	http.Error(r.w, s, code)
}

// LoadFile handling file uploads
func (r *Request) LoadFile(field, dir string) (string, error) {
	r.req.ParseMultipartForm(32 << 20)
	file, handler, err := r.req.FormFile(field)
	if err != nil {
		return "", err
	}
	defer file.Close()
	fmt.Fprintf(r.w, "%v", handler.Header)
	f, err := os.OpenFile(dir+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()
	io.Copy(f, file)
	return handler.Filename, nil
}

type parts struct {
	id, action string
}

func (p parts) ID() (i int64) {
	i, _ = strconv.ParseInt(p.id, 10, 64)
	return
}

func getParts(path, prefix string) parts {
	res := parts{}
	path = strings.TrimPrefix(path, prefix)
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	parts := strings.Split(path, "/")

	l := len(parts)

	if l > 0 {
		res.id = parts[0]
		if l > 1 {
			res.action = parts[1]
		}
	}

	return res
}
