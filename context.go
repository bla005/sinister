package sinister

import (
	"encoding/json"
	"net"
	"net/http"

	"go.uber.org/zap"
)

// HC ...
type HC struct {
	w      http.ResponseWriter
	r      *http.Request
	logger *zap.Logger
	params []*Param
}

func (hc *HC) set(w http.ResponseWriter, r *http.Request, logger *zap.Logger, params []*Param) {
	hc.w = w
	hc.r = r
	hc.logger = logger
	hc.params = params
}

func (hc *HC) reset() {
	hc.w = nil
	hc.r = nil
	hc.logger = nil
	hc.params = nil
}

func newHC() *HC {
	return &HC{
		w:      nil,
		r:      nil,
		logger: nil,
		params: nil,
	}
}

// JSONS ...
func (hc *HC) JSONS(code int, data string) {
	// r := &httpResponse{code, msg}
	r := newHTTPResponse(code, data)
	hc.w.WriteHeader(code)
	if err := json.NewEncoder(hc.w).Encode(r); err != nil {
		http.Error(hc.w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// JSONI ...
func (hc *HC) JSONI(code int, data interface{}) {
	hc.w.WriteHeader(code)
	if err := json.NewEncoder(hc.w).Encode(data); err != nil {
		http.Error(hc.w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// RAWS ...
func (hc *HC) RAWS(code int, data string) {
	_, err := hc.w.Write([]byte(data))
	if err != nil {
		http.Error(hc.w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// RAWB ...
func (hc *HC) RAWB(code int, data []byte) {
	_, err := hc.w.Write(data)
	if err != nil {
		http.Error(hc.w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// Param ...
func (hc *HC) Param(param string) URLParam {
	if len(hc.params) == 0 {
		return ""
	}
	return URLParam(findParam(hc.params, param))
}

// Log ...
func (hc *HC) Log(msg string, level LogLevel) {
	switch level {
	case DEBUG:
		hc.logger.Debug(msg)
	case ERROR:
		hc.logger.Error(msg)
	case INFO:
		hc.logger.Info(msg)
	case FATAL:
		hc.logger.Fatal(msg)
	case WARN:
		hc.logger.Warn(msg)
	}
}

// Query ...
func (hc *HC) Query(key string) (string, error) {
	if hc.r.URL.Query().Get(key) == "" {
		return "", ErrQueryNotFound
	}
	return hc.r.URL.Query().Get(key), nil
}

// MIME ...
func (hc *HC) MIME(mime MIME) {
	hc.w.Header().Set("Content-Type", string(mime))
}

// ClientIP ...
func (hc *HC) ClientIP() string {
	ip, _, err := net.SplitHostPort(hc.r.RemoteAddr)
	if err != nil {
		forward := hc.r.Header.Get("X-Forwarded-For")
		return forward
	}
	return ip
}
