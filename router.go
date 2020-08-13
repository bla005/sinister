package sinister

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"go.uber.org/zap"
)

// ErrInvalidParam ...
var ErrInvalidParam = errors.New("sinister: invalid param")

// Handler ...
type Handler func(*HC)

type router struct {
	routes []*route
	//prefixes []*prefix
	node            *node
	pool            *sync.Pool
	NotFoundHandler Handler
}

/*
type prefix struct {
	prefix string
	routes []*route
}

func newPrefix(pref string) *prefix {
	return &prefix{
		prefix: pref,
		routes: make([]*route, 0),
	}
}
*/
func newRouter() *router {
	return &router{
		routes: nil,
		//prefixes: make([]*prefix, 0),
		node: nil,
		pool: &sync.Pool{
			New: func() interface{} { return newHC() },
		},
	}
}

func findParam(params []*Param, param string) string {
	for _, p := range params {
		if p.Name == param {
			return p.Value
		}
	}
	return ""
}

type route struct {
	path    string
	rawPath string
	method  string
	handler Handler
	params  []string
}

func newRoute(path, rawPath, method string, h Handler, params []string) *route {
	return &route{
		path:    path,
		rawPath: rawPath,
		method:  method,
		handler: h,
		params:  params,
	}
}

func setParams(params []string, values []string) []*Param {
	if len(params) == 0 || len(params) != len(values) {
		return nil
	}
	paramsOut := make([]*Param, len(params))
	param := &Param{}
	for i, v := range params {
		param = &Param{Name: v, Value: values[i]}
		paramsOut[i] = param
	}
	return paramsOut
}

// ServeHTTP ...
func (router *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	formattedPath, params, valid := validateRequestPath(r.URL.Path, r.Method)
	fmt.Println("serve")
	fmt.Println(formattedPath, params, valid)
	if valid {
		route := findNode(router.node, formattedPath)
		ctxLogger, _ := zap.NewProduction()
		lib := router.pool.Get().(*HC)
		lib.reset()
		if route != nil && isMatch(route.rawPath, formattedPath) {
			fmt.Println("is match")
			urlParams := setParams(route.params, params)
			lib.set(w, r, ctxLogger, urlParams)
			route.handler(lib)
		} else {
			lib.set(w, r, ctxLogger, nil)
			router.NotFoundHandler(lib)
		}
		router.pool.Put(lib)
		return
	}
}
