package rapi

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"path"
	"strings"
	"unicode"
)

func extractJSONPayload(data io.Reader, v interface{}) error {
	decoder := json.NewDecoder(data)
	err := decoder.Decode(&v)
	return err
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	a := []rune(s)
	a[0] = unicode.ToUpper(a[0])
	return string(a)
}

func urlParts(path string) map[int]string {
	res := make(map[int]string)
	path = strings.TrimSpace(path)
	path = strings.TrimLeft(path, "/")
	path = strings.TrimRight(path, "/")
	parts := strings.Split(path, "/")
	for k, v := range parts {
		res[k] = v
	}
	return res
}

// RenderJSONError common function to render error to client in JSON format
func RenderJSONError(w http.ResponseWriter, code int, s string) {
	RenderJSON(w, code, JSONData{"errors": JSONData{"message": []string{s}}})
}

// RenderJSON common function to render JSON to client
func RenderJSON(w http.ResponseWriter, code int, s JSONData) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(s)
}

// RenderJSONgzip common function to render gzipped JSON to client
func RenderJSONgzip(w http.ResponseWriter, code int, s JSONData) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(code)
	gz := gzip.NewWriter(w)
	defer gz.Close()
	json.NewEncoder(gz).Encode(s)
}

// cleanPath returns the canonical path for p, eliminating . and .. elements.
// Borrowed from the net/http package.
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}
	return np
}
