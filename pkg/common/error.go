package common

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/pkg/errors"

	"github.com/NexClipper/logger"
)

// StandardError standard error for klevr
type StandardError struct {
	message   string
	timestamp int64
	cause     string
}

// HTTPError runtime error for klevr
type HTTPError struct {
	StandardError
	statusCode int
}

func (e *StandardError) Error() string {
	return fmt.Sprintf("StandardError : [\nmessage : %s\ntimestamp : %d\ncause : %s]]",
		e.message, e.timestamp, e.cause)
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTPError : [\nstatusCode : %d\nmessage : %s\ntimestamp : %d\ncause : %s]]",
		e.statusCode, e.message, e.timestamp, e.cause)
}

func (e *StandardError) initStandardError(message string, err error) {
	if err != nil {
		e.message = message + "\nroot message : " + err.Error()
	} else {
		e.message = message
	}

	e.timestamp = time.Now().Unix()
	e.cause = string(debug.Stack())
}

// NewStandardError constructor for StandardError
func NewStandardError(message string) *StandardError {
	e := &StandardError{
		timestamp: time.Now().UTC().Unix(),
	}

	e.initStandardError(message, nil)

	return e
}

// NewStandardErrorWrap constructor for StandardError
func NewStandardErrorWrap(message string, err error) *StandardError {
	e := &StandardError{
		timestamp: time.Now().UTC().Unix(),
	}

	e.initStandardError(message, err)

	return e
}

func (e *HTTPError) initHTTPError(statusCode int, message string, err error) {
	e.initStandardError(message, err)

	e.statusCode = statusCode
}

// NewHTTPError constructor for HTTPError
func NewHTTPError(statusCode int, message string) *HTTPError {
	e := &HTTPError{}

	e.initHTTPError(statusCode, message, nil)

	return e
}

// NewHTTPErrorWrap constructor for HTTPError
func NewHTTPErrorWrap(statusCode int, message string, err error) *HTTPError {
	e := &HTTPError{}

	e.initHTTPError(statusCode, message, err)

	return e
}

// ErrorWithPanic raise panic with RuntimeError when error is not nil.
func ErrorWithPanic(err error, message string) {
	if err != nil {
		panic(errors.Wrap(err, message))
	}
}

// ErrorWithDebugLog log with specified log level
func ErrorWithDebugLog(err error, message string) {
	errorWithLog(err, 0, message)
}

// ErrorWithInfoLog log with specified log level
func ErrorWithInfoLog(err error, message string) {
	errorWithLog(err, 1, message)
}

// ErrorWithWarnLog log with specified log level
func ErrorWithWarnLog(err error, message string) {
	errorWithLog(err, 2, message)
}

// ErrorWithErrorLog log with specified log level
func ErrorWithErrorLog(err error, message string) {
	errorWithLog(err, 3, message)
}

// log with specified log level
func errorWithLog(err error, l logger.Level, message string) {
	const logFormat string = "%s : %+v"
	switch l {
	case 0:
		logger.Debugf(logFormat, message, err)
	case 1:
		logger.Infof(logFormat, message, err)
	case 2:
		logger.Warningf(logFormat, message, err)
	case 3:
		logger.Errorf(logFormat, message, err)
	case 4:
		logger.Fatalf(logFormat, message, err)
	default:
		logger.Debugf(logFormat, message, err)
	}
}

// Block {
// 	Try: func() {
// 		fmt.Println("Try..")
// 		Throw("stop it")
// 	},
// 	Catch: func(e Exception) {
// 		fmt.Printf("Caught %v\n", e)
// 	},
// 	Finally: func() {
// 		fmt.Println("Finally..")
// 	},
// }.Do()
// Block Try-Catch-Finally block struct
type Block struct {
	Try     func()
	Catch   func(error)
	Finally func()
}

// Exception pass exception to Catch
type Exception interface{}

// Throw raise panic with exception
func Throw(up Exception) {
	panic(up)
}

// Do run block state
func (b Block) Do() {
	if b.Finally != nil {
		defer b.Finally()
	}

	if b.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				b.Catch(errors.WithStack(fmt.Errorf("%+v", r)))
			}
		}()
	}

	b.Try()
}
