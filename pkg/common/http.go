package common

import (
	"net/http"
	"strconv"
	"strings"
)

// Request http.Request override
type Request struct {
	*http.Request
}

// ResponseWrapper http.ResponseWriter override
type ResponseWrapper struct {
	http.ResponseWriter

	StatusCode int
}

// Param return string param
func (r *Request) Param(name string) string {
	return strings.Join(r.URL.Query()[name], "")
}

// ParamToInt return int param
func (r *Request) ParamToInt(name string) (int, error) {
	return strconv.Atoi(r.Param(name))
}

// ParamToBool return bool param
func (r *Request) ParamToBool(name string) (bool, error) {
	return strconv.ParseBool(r.Param(name))
}

// ParamToUInt return int param
func (r *Request) ParamToUInt(name string) (uint64, error) {
	return strconv.ParseUint(r.Param(name), 10, 64)
}

// Header ResponseWriter Header override
func (w *ResponseWrapper) Header() http.Header {
	return w.ResponseWriter.Header()
}

// Write ResponseWriter Write override
func (w *ResponseWrapper) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

// WriteHeader ResponseWriter WriteHeader override
func (w *ResponseWrapper) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
