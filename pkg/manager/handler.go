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
)

// CommonWrappingHandler common handler for processing standard
func CommonWrappingHandler(ctx *common.Context) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var tx *Tx

			// 복잡한 flow 정리를 위해 try-catch-finally 블럭 사용
			common.Block{
				Try: func() {
					// response wrapping
					nw := common.ResponseWrapper{
						ResponseWriter: w,
						StatusCode:     http.StatusOK,
					}

					w.Header().Set("Content-Type", "json/application; charset=utf-8")
					w.Header().Set("X-Content-Type-Options", "nosniff")

					// DB session 시작
					tx = &Tx{CtxGetDbConn(ctx).NewSession()}
					err := tx.Begin()
					if err != nil {
						logger.Errorf("DB session begin error : %v", err)
						common.Throw(err)
					}

					// request context 생성
					rCtx := common.FromContext(ctx)
					// Request context에 DB session 설정
					rCtx.Put(CtxDbSession, tx)

					// request에 context 추가
					context.Set(r, CtxRequestContext, rCtx)

					// 다음 핸들러로 진행
					next.ServeHTTP(&nw, r)

					// 트랜잭션 commit
					err = tx.Commit()
					if err != nil {
						logger.Warningf("commit failed : %v", err)
						common.Throw(err)
					}
				},
				Catch: func(e common.Exception) {
					// 트랜잭션 recover 정의
					if !tx.IsClosed() {
						tx.Rollback()
					}

					common.WriteHTTPError(500, w, fmt.Errorf("%+v", e), "Service is unavailable")
				},
				Finally: func() {
					// Context 초기화
					defer func() {
						context.Clear(r)
					}()

					// 세션 close
					defer func() {
						if !tx.IsClosed() {
							tx.Close()
						}
					}()
				},
			}.Do()
		})
	}
}

// RequestInfoLoggerHandler request processing information logging handler
func RequestInfoLoggerHandler(next http.Handler) http.Handler {
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
