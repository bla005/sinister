package sinister

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"go.uber.org/zap"
)

type Handler func(*Lib)

func findParam(params []*Param, param string) string {
	for _, p := range params {
		if p.Name == param {
			return p.Value
		}
	}
	return ""
}

type Lib struct {
	w      http.ResponseWriter
	r      *http.Request
	logger *zap.Logger
	params []*Param
}

func NewLib(w http.ResponseWriter, r *http.Request, logger *zap.Logger) *Lib {
	return &Lib{
		w:      w,
		r:      r,
		logger: nil,
		params: make([]*Param, 0),
	}
}

type httpResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (l *Lib) JSON(code int, msg string) {
	r := &httpResponse{code, msg}
	if err := json.NewEncoder(l.w).Encode(r); err != nil {
		http.Error(l.w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
func (l *Lib) set(w http.ResponseWriter, r *http.Request, params []*Param) {
	l.w = w
	l.r = r
	l.params = params
}

func (l *Lib) reset() {
	l.w = nil
	l.r = nil
	l.params = nil
}
func (l *Lib) Query(key string) string {
	return l.r.URL.Query().Get(key)
}

type ParamValue string

func (p ParamValue) Int() (*int, error) {
	n, err := strconv.Atoi(string(p))
	if err != nil {
		return nil, ErrInvalidParam
	}
	return n, nil
}
func (p ParamValue) Int64() (*int64, error) {
	n, err := strconv.Atoi(string(p))
	if err != nil {
		return nil, ErrInvalidParam
	}
	return int64(n), nil
}
func (p ParamValue) String() string {
	return string(p)
}
func (p ParamValue) Bytes() []byte {
	return []byte(p)
}

var (
	ErrInvalidParam = errors.New("lol")
)

func (l *Lib) Param(param string) ParamValue {
	if len(l.params) == 0 {
		return ""
	}
	return ParamValue(findParam(l.params, param))
}

type param struct {
	name string
	pos  int
}
type Route struct {
	Name    string
	Path    string
	RawPath string
	Method  string
	Handler Handler
	params  []string
}
type Router struct {
	Routes []*Route
	Node   *node
	Pool   *sync.Pool
}

func NewRouter() *Router {
	return &Router{
		Routes: make([]*Route, 0),
		Node:   nil,
		Pool: &sync.Pool{
			New: func() interface{} { return &Lib{} },
		},
	}
}

func (r *Router) Get(name, path string, h Handler) {
	params, fPath := validatePath(path)
	fmt.Println("validatePath", params, fPath)
	r1 := NewRoute(name, path, fPath, http.MethodGet, h, params)
	r.Node = insert(r.Node, r1)
}

func (r *Route) Param(name string) int {
	// return findParam(r.params, name)
	return 0
}

type Param struct {
	Name  string
	Value string
}

func NewRoute(name, path, rawPath, method string, h Handler, params []string) *Route {
	return &Route{
		Name:    name,
		Path:    path,
		RawPath: rawPath,
		Method:  method,
		Handler: h,
		params:  params,
	}
}

func setParams(params []string, values []string) []*Param {
	if len(params) == 0 || len(params) != len(values) {
		return nil
	}
	rp := make([]*Param, len(params))
	t := &Param{}
	for i, p := range params {
		t = &Param{Name: p, Value: values[i]}
		rp[i] = t
	}
	return rp
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request path", r.URL.Path)
	if isReqPathValid(r.URL.Path) {
		fmt.Println("valid request url", r.URL.Path)
		fPath, pParams := formatRequestPath(r.URL.Path)
		route := findNode(router.Node, fPath)
		fmt.Println("findNode", route)
		if route != nil && isMatch(route.RawPath, fPath) {
			ep := setParams(route.params, pParams)
			lib := router.Pool.Get().(*Lib)
			lib.set(w, r, ep)
			route.Handler(lib)
			lib.reset()
			router.Pool.Put(lib)
			return
		}
	}
	http.NotFound(w, r)
}
