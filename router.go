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
	logger          *zap.Logger
}

func (r *router) setNotFoundHandler() {
	r.NotFoundHandler = func(ctx *HC) {
		ctx.MIME(ApplicationJSON)
		ctx.JSONS(404, "Not found")
	}
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
func newRouter(logger *zap.Logger) *router {
	return &router{
		routes: nil,
		//prefixes: make([]*prefix, 0),
		node: nil,
		pool: &sync.Pool{
			New: func() interface{} { return newHC() },
		},
		logger: logger,
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
		// ctxLogger, _ := zap.NewProduction()
		nLog := router.logger
		lib := router.pool.Get().(*HC)
		lib.reset()
		if route != nil && isMatch(route.rawPath, formattedPath) {
			fmt.Println("is match")
			urlParams := setParams(route.params, params)
			lib.set(w, r, nLog, urlParams)
			route.handler(lib)
		}
		lib.set(w, r, nLog, nil)
		router.NotFoundHandler(lib)
		router.pool.Put(lib)
		fmt.Println("ok")
	}
}
