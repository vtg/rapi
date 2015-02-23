rapi
====
HTTP routing package that helps to create restfull json api for Go applications.

Note: The project has been moved into new repository: [github.com/vtg/flash](https://github.com/vtg/flash)

what it does:

 - dispatching actions to controllers
 - rendering JSON response
 - extracting JSON request data by key
 - handling file uploads
 - sending gzipped JSON responses when applicable
 - sending gzipped versions of static files if any

standard REST usage example:
   
```go
package main

import (
	"net/http"

	"github.com/vtg/rapi"
)

var pages map[int64]*Page

func main() {
	pages = make(map[int64]*Page)
	pages[1] = &Page{Id: 1, Name: "Page 1"}
	pages[2] = &Page{Id: 2, Name: "Page 2"}
	r := rapi.NewRouter()
	a := r.PathPrefix("/api/v1")
	// see Route.Route for more info
	a.Route("/pages", &Pages{}, "page", authenticate)
	// see Route.FileServer for more info
	r.PathPrefix("/images/").FileServer("./public/")
	r.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8090", r)
}

// simple quthentication implementation
func authenticate(c rapi.Controller) bool {
	key := c.QueryParam("key")
	if key == "pass" {
		return true
	} else {
		c.RenderJSONError(http.StatusUnauthorized, "unauthorized")
	}
	return false
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

type Page struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

func findPage(id int64) *Page {
	p := pages[id]
	return p
}
func insertPage(p Page) *Page {
	id := int64(len(pages) + 1)
	p.Id = id
	pages[id] = &p
	return pages[id]
}

// Pages used as controller
type Pages struct {
	rapi.Request
}

// Index processed on GET /pages
func (p *Pages) Index() {
	pgs := []Page{}
	for _, v := range pages {
		pgs = append(pgs, *v)
	}
	p.RenderJSON(200, rapi.JSONData{"pages": pgs})
}

// Show processed on GET /pages/1
func (p *Pages) Show() {
	page := findPage(p.ID)
	if page.Id > 0 {
		p.RenderJSON(200, rapi.JSONData{"page": page})
	} else {
		p.RenderJSONError(404, "record not found")
	}
}

// Create processed on POST /pages
// with input data provided {"page":{"name":"New Page","content":"some content"}}
func (p *Pages) Create() {
	m := Page{}
	// see Request.LoadJSONRequest for more info
	p.LoadJSONRequest(p.Root, &m)
	if m.Name == "" {
		p.RenderJSONError(422, "name required")
	} else {
		insertPage(m)
		p.RenderJSON(200, rapi.JSONData{p.Root: m})
	}
}

// Update processed on PUT /pages/1
// with input data provided {"page":{"name":"Page 1","content":"updated content"}}
func (p *Pages) Update() {
	m := Page{}
	p.LoadJSONRequest(p.Root, &m)
	page := findPage(p.ID)
	if page.Id > 0 {
		page.Content = m.Content
		p.RenderJSON(200, rapi.JSONData{"page": page})
	} else {
		p.RenderJSONError(404, "record not found")
	}
}

// Destroy processed on DELETE /pages/1
func (p *Pages) Destroy() {
	page := findPage(p.ID)
	if page.Id > 0 {
		delete(pages, page.Id)
		p.RenderJSON(203, rapi.JSONData{})
	} else {
		p.RenderJSONError(404, "record not found")
	}
}
```

Its possible to serve custom actions.  
To add custom action to controller prefix action name with HTTP method:

```go
 // POST /pages/clean or POST /pages/1/clean
 func (p *Pages) POSTClean {
   // do some work here
 }
 // DELETE /pages/clean or DELETE /pages/1/clean
 func (p *Pages) DELETEClean {
   // do some work here
 }
 // GET /pages/stat or GET /pages/1/stat
 func (p *Pages) GETStat {
   // do some work here
 }
 ...
```

#####Author

VTG - http://github.com/vtg

##### License

Released under the [MIT License](http://www.opensource.org/licenses/MIT).

[![GoDoc](https://godoc.org/github.com/vtg/rapi?status.png)](http://godoc.org/github.com/vtg/rapi)
