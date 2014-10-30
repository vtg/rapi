package rapi

import (
	"net/http"
	"os"
	"path"
	"strings"
)

type fileHandler struct {
	root         string
	dirList, enc bool
}

func (f *fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !f.dirList {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
	}

	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	upath = path.Join(f.root, upath)
	upath = path.Clean(upath)

	if f.enc {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			p := path.Ext(upath)
			switch p {
			case ".css", ".js":
				if _, err := os.Stat(upath + ".gz"); err == nil {
					upath = upath + ".gz"
					if p == ".js" {
						w.Header().Set("Content-Type", "application/javascript")
					} else {
						w.Header().Set("Content-Type", "text/css; charset=utf-8")
					}
					w.Header().Set("Content-Encoding", "gzip")
				}
			}
		}
	}

	http.ServeFile(w, r, upath)
}

// fileServer returns a handler that serves HTTP requests
// with the contents of the file system rooted at root.
//
// first boolean value is dirListing enable/disable. default false
//
// second boolean value is serving gzip encoded files. default false
func fileServer(root string, bools []bool) http.Handler {
	f := &fileHandler{root: root}
	if len(bools) > 0 {
		f.dirList = bools[0]
	}
	if len(bools) > 1 {
		f.enc = bools[1]
	}
	return f
}
