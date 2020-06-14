package sinister

import (
	"fmt"
	"net/http"
)

type Sinister struct {
	router      *router
	Middlewares []*Middleware
}
type Middleware func(*Lib) *Lib

func (s *Sinister) Get(path string, h Handler) {
	params, formattedPath, encoded := validatePath(path)
	fmt.Println("validatePath", params, formattedPath)
	r1 := newRoute(path, formattedPath, http.MethodGet, h, params, encoded)
	s.router.Node = insert(s.router.Node, r1)
}

func New() *Sinister {
	return &Sinister{
		router:      newRouter(),
		Middlewares: make([]*Middleware, 0),
	}
}

func (lib *Sinister) UseMiddleware(m *Middleware) {
	lib.Middlewares = append(lib.Middlewares, m)
}
