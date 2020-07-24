package manager

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/gorilla/context"

	"github.com/gorilla/mux"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"xorm.io/xorm"
)

// CommonWrappingHandler common handler for processing standard
func CommonWrappingHandler(DB *xorm.Engine) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// response wrapping
			nw := common.ResponseWrapper{w, http.StatusOK}

			// DB session 시작
			session := DB.NewSession()

			err := session.Begin()
			if err != nil {
				logger.Errorf("DB session begin error : %v", err)
				common.HTTPError(500, w, err, "Service is unavailable")
			}

			// 트랜잭션 recover 정의
			defer func() {
				r := recover()
				if r != nil {
					logger.Warningf("recovered : %v", r)

					if !session.IsClosed() {
						session.Rollback()
					}

					common.HTTPError(500, w, common.NewRuntimeError(fmt.Sprintf("%", r)), "Service is unavailable")
				}
			}()

			// Request context에 DB session 설정
			context.Set(r, DBConnContextName, session)

			// 다음 핸들러로 진행
			next.ServeHTTP(&nw, r)

			// 트랜잭션 commit
			err = session.Commit()
			if err != nil {
				logger.Warningf("commit failed : %v", err)

				common.HTTPError(500, w, common.NewRuntimeError(fmt.Sprintf("%", err)), "Service is unavailable")
			}

			// 세션 close
			defer func() {
				if !session.IsClosed() {
					session.Close()
				}
			}()

		})
	}
}

// // CommonWrappingHandler common handler for processing standard
// func CommonWrappingHandler(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		// response wrapping
// 		nw := common.ResponseWrapper{w, http.StatusOK}

// 		next.ServeHTTP(&nw, r)
// 	})
// }

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

		// Process request
		next.ServeHTTP(w, r)

		nw := reflect.ValueOf(w).Interface().(*common.ResponseWrapper)

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
		param.StatusCode = nw.StatusCode
		// param.ErrorMessage = nw.Result().c.Errors.ByType(ErrorTypePrivate).String()

		// param.BodySize = c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		param.Path = path

		logger.Debug(formatter(param))
	})
}
