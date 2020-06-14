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
		logger: logger,
		params: make([]*Param, 0),
	}
}

type httpResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (l *Lib) String(code int, msg string) {
	_, err := l.w.Write([]byte(msg))
	if err != nil {
		http.Error(l.w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
func (l *Lib) Log(msg string) {
}

func (l *Lib) JSON(code int, msg string) {
	r := &httpResponse{code, msg}
	l.w.WriteHeader(code)
	if err := json.NewEncoder(l.w).Encode(r); err != nil {
		http.Error(l.w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
func (l *Lib) JSONv2(code int, data interface{}) {
	l.w.WriteHeader(code)
	if err := json.NewEncoder(l.w).Encode(data); err != nil {
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

func (p ParamValue) Int() (int, error) {
	n, err := strconv.Atoi(string(p))
	if err != nil {
		return 0, ErrInvalidParam
	}
	return n, nil
}
func (p ParamValue) Int64() (int64, error) {
	n, err := strconv.Atoi(string(p))
	if err != nil {
		return 0, ErrInvalidParam
	}
	m := int64(n)
	return m, nil
}
func (p ParamValue) String() string {
	return string(p)
}
func (p ParamValue) Bytes() []byte {
	return []byte(p)
}

var (
	ErrInvalidParam = errors.New("sinister: invalid param")
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
type route struct {
	path    string
	rawPath string
	method  string
	handler Handler
	params  []string
	encoded int
}

func newRoute(path, rawPath, method string, h Handler, params []string, encoded int) *route {
	return &route{
		path:    path,
		rawPath: rawPath,
		method:  method,
		handler: h,
		params:  params,
		encoded: encoded,
	}
}

type router struct {
	Routes []*route
	Node   *node
	Pool   *sync.Pool
}

func newRouter() *router {
	return &router{
		Routes: make([]*route, 0),
		Node:   nil,
		Pool: &sync.Pool{
			New: func() interface{} { return &Lib{} },
		},
	}
}

func (r *router) Get(name, path string, h Handler) {
	params, formattedPath, encoded := validatePath(path)
	fmt.Println("validatePath", params, formattedPath)
	r1 := newRoute(path, formattedPath, http.MethodGet, h, params, encoded)
	r.Node = insert(r.Node, r1)
}

func (r *route) Param(name string) int {
	// return findParam(r.params, name)
	return 0
}

type Param struct {
	Name  string
	Value string
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

func (router *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request path", r.URL.Path)
	formattedPath, params, encoded, valid := formatReqPath(r.URL.Path)
	if valid {
		fmt.Println("valid request url", r.URL.Path)
		// fPath, pParams := formatRequestPath(r.URL.Path)
		// route := findNode(router.Node, formattedPath)
		route := findNode(router.Node, encoded)
		fmt.Println("findNode", route)
		if route != nil && isMatch(route.RawPath, formattedPath) {
			ep := setParams(route.params, params)
			lib := router.Pool.Get().(*Lib)
			// logger := lib.logger.With(zap.String("path", r.URL.Path))
			lib.set(w, r, ep)
			route.Handler(lib)
			lib.reset()
			router.Pool.Put(lib)
			return
		}
	}
	http.NotFound(w, r)
}
