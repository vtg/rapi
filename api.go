package rapi

import (
	"net/http"
	"reflect"
)

type Controller interface {
	Init(w http.ResponseWriter, req *http.Request, root, prefix string, extras []string)
	QueryParam(string) string
	SetParam(string, interface{})
	Param(string) interface{}
	Header(string) string
	CurrentAction() string
	RenderJSON(code int, s JSONData)
	RenderJSONError(code int, s string)
}

// ReqFunc is the function type for middlware
type ReqFunc func(Controller) bool

// JSONData shortcut for map[string]interface{}
type JSONData map[string]interface{}

// handle returns http handler function that will process controller actions
func handle(i Controller, rootKey, prefix string, extras []string, funcs ...ReqFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		t := reflect.Indirect(reflect.ValueOf(i)).Type()
		c := reflect.New(t)
		ctr := c.Interface().(Controller)
		ctr.Init(w, req, rootKey, prefix, extras)

		for _, f := range funcs {
			if ok := f(ctr); !ok {
				return
			}
		}

		if method := c.MethodByName(ctr.CurrentAction()); method.IsValid() {
			method.Call([]reflect.Value{})
		} else {
			RenderJSONError(w, http.StatusBadRequest, "action not found")
		}
	}
}
