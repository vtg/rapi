package rapi

import (
	"net/http"
	"reflect"
)

type Controller interface {
	Init(w http.ResponseWriter, req *http.Request, root string)
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
func handle(i Controller, rootKey string, funcs []ReqFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		t := reflect.Indirect(reflect.ValueOf(i)).Type()
		c := reflect.New(t)
		c.MethodByName("Init").Call([]reflect.Value{
			reflect.ValueOf(w),
			reflect.ValueOf(req),
			reflect.ValueOf(rootKey),
		})

		for _, f := range funcs {
			if ok := f(c.Interface().(Controller)); !ok {
				return
			}
		}

		action := c.MethodByName("CurrentAction").Call([]reflect.Value{})

		if method := c.MethodByName(action[0].String()); method.IsValid() {
			method.Call([]reflect.Value{})
		} else {
			RenderJSONError(w, http.StatusBadRequest, "action not found")
		}
	}
}
