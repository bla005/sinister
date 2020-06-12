package sinister

type Sinister struct {
	Router      *Router
	Middlewares []*Middleware
}

func New(router *Router) *Sinister {
	return &Sinister{
		Router:      router,
		Middlewares: make([]*Middleware, 0),
	}
}

type Middleware func(*Lib) *Lib

func (lib *Sinister) UseMiddleware(m *Middleware) {
	lib.Middlewares = append(lib.Middlewares, m)
}
