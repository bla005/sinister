package sinister

import (
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const (
	// GET ...
	GET = http.MethodGet
	// POST ...
	POST = http.MethodPost
	// PUT ...
	PUT = http.MethodPut
	// DELETE ...
	DELETE = http.MethodDelete
)

var (
	// ErrQueryNotFound ...
	ErrQueryNotFound error
)

const (
	charsetUTF8 = "charset=UTF-8"
)

// MIME ...
type MIME string

const (
	// ApplicationJSON ...
	ApplicationJSON MIME = "application/json"
	// ApplicationJavascript ...
	ApplicationJavascript MIME = "application/javascript"
	// ApplicationXML ...
	ApplicationXML MIME = "application/xml"
	// TextXML ...
	TextXML MIME = "text/xml"
	// ApplicationForm ...
	ApplicationForm MIME = "application/x-www-form-urlencoded"
	// ApplicationProtobuf ...
	ApplicationProtobuf MIME = "application/protobuf"
	// ApplicationMsgpack ...
	ApplicationMsgpack MIME = "application/msgpack"
	// TextHTML ...
	TextHTML MIME = "text/html"
	// TextPlain ...
	TextPlain MIME = "text/plain"
	// MultipartForm ...
	MultipartForm MIME = "multipart/form-data"
	// OctetStream ...
	OctetStream MIME = "application/octet-stream"
)

const (
	// ApplicationJSONCharsetUTF8 ...
	ApplicationJSONCharsetUTF8 MIME = ApplicationJSON + "; " + charsetUTF8
	// ApplicationJavascriptCharsetUTF8 ...
	ApplicationJavascriptCharsetUTF8 MIME = ApplicationJavascript + "; " + charsetUTF8
	// ApplicationXMLCharsetUTF8 ...
	ApplicationXMLCharsetUTF8 MIME = ApplicationXML + "; " + charsetUTF8
	// TextXMLCharsetUTF8 ...
	TextXMLCharsetUTF8 MIME = TextXML + "; " + charsetUTF8
	// TextHTMLCharsetUTF8 ...
	TextHTMLCharsetUTF8 MIME = TextHTML + "; " + charsetUTF8
	// TextPlainCharsetUTF8 ...
	TextPlainCharsetUTF8 MIME = TextPlain + "; " + charsetUTF8
)

// Sinister ...
type Sinister struct {
	logger      *zap.Logger
	router      *router
	middlewares []*Middleware
	server      *http.Server
}

// Middleware ...
type Middleware func(*HC) *HC

func (s *Sinister) register(path, method string, h Handler) {
	params, formattedPath := validatePath(path, method)
	r1 := newRoute(path, formattedPath, method, h, params)
	s.router.node = insert(s.router.node, r1)
}

// GET ...
func (s *Sinister) GET(path string, h Handler) {
	s.register(path, GET, h)
}

// POST ...
func (s *Sinister) POST(path string, h Handler) {
	s.register(path, POST, h)
}

// PUT ...
func (s *Sinister) PUT(path string, h Handler) {
	s.register(path, PUT, h)
}

// DELETE ...
func (s *Sinister) DELETE(path string, h Handler) {
	s.register(path, DELETE, h)
}

// New ...
func New() *Sinister {
	return &Sinister{
		logger:      newLogger(),
		router:      newRouter(),
		middlewares: make([]*Middleware, 0),
		server:      &http.Server{},
	}
}

// Start ...
func (s *Sinister) Start(addr string) error {
	// if err := http.ListenAndServe(":8080", s.router); err != nil {
	// log.Fatal(err)
	// }
	s.server.Handler = s.router
	l, err := newListener(addr)
	if err != nil {
		return err
	}
	return s.server.Serve(l)
}

// Close ...
func (s *Sinister) Close() error {
	return s.server.Close()
}

// UseMiddleware ...
func (s *Sinister) UseMiddleware(m *Middleware) {
	s.middlewares = append(s.middlewares, m)
}

// HTTPResponse ...
type HTTPResponse struct {
	Code    int         `json:"code"`
	Message interface{} `json:"msg"`
}

func newHTTPResponse(code int, msg interface{}) *HTTPResponse {
	return &HTTPResponse{
		Code:    code,
		Message: msg,
	}
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	c, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	if err := c.SetKeepAlive(true); err != nil {
		return nil, err
	}
	if err := c.SetKeepAlivePeriod(1 * time.Minute); err != nil {
		return nil, err
	}
	return c, nil

	/*
		if c, err = ln.AcceptTCP(); err != nil {
			return
		} else if err = c.(*net.TCPConn).SetKeepAlive(true); err != nil {
			return
		}
		_ = c.(*net.TCPConn).SetKeepAlivePeriod(3 * time.Minute)
		return
	*/
}

func newListener(addr string) (*tcpKeepAliveListener, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &tcpKeepAliveListener{listener.(*net.TCPListener)}, nil
}
