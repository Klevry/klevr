package common

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/NexClipper/logger"
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

// WriteHTTPError common error process
func WriteHTTPError(statusCode int, w http.ResponseWriter, err error, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "%s : %v", message, err)

	logger.Warningf("%s : %+v\n%s", message, err, debug.Stack())
}
