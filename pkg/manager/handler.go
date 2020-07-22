package manager

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"github.com/gin-gonic/gin"
)

// APIHandler API handler
func (api *API) APIHandler(h func(*gin.Context)) gin.HandlerFunc {
	logger.Debug("pre handle")

	return h
}

// TestHandler test handler
func TestHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("호출전")

		// Process request
		next.ServeHTTP(w, r)

		logger.Debug("호출후")
	})
}

// ExecutionInfoLoggerHandler request processing information logging handler
func ExecutionInfoLoggerHandler(next http.Handler) http.Handler {
	var formatter = func(param common.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		if param.IsOutputColor() {
			statusColor = param.StatusCodeColor()
			methodColor = param.MethodColor()
			resetColor = param.ResetColor()
		}

		if param.Latency > time.Minute {
			// Truncate in a golang < 1.8 safe way
			param.Latency = param.Latency - param.Latency%time.Second
		}
		return fmt.Sprintf("|%s %3d %s| %9v | %15s |%s %-7s %s %#v\n%s",
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			param.Path,
			param.ErrorMessage,
		)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Start timer
		start := time.Now()
		path := r.URL.Path
		raw := r.URL.RawQuery

		nw := httptest.NewRecorder()

		// Process request
		next.ServeHTTP(nw, r)

		keys := make(map[string]interface{})
		for k, v := range r.URL.Query() {
			keys[k] = v
		}

		param := common.LogFormatterParams{
			Request: r,
			Keys:    keys,
		}

		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)
		param.ClientIP = r.RemoteAddr
		param.Method = r.Method
		param.StatusCode = nw.Result().StatusCode
		// param.ErrorMessage = nw.Result().c.Errors.ByType(ErrorTypePrivate).String()

		// param.BodySize = c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		param.Path = path

		logger.Debug(formatter(param))
	})
}
