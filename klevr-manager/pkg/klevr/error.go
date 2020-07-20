package klevr

import (
	"fmt"
	"time"
)

// StandardError standard error for klevr
type StandardError struct {
	error
	code      int
	message   string
	timestamp int64
	cause     error
}

// RuntimeError runtime error for klevr
type RuntimeError struct {
	StandardError
}

// CheckedError checked error for klevr
type CheckedError struct {
	StandardError
}

func (e *StandardError) Error() string {
	return fmt.Sprintf("StandardError : %v", e)
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("RuntimeError : %v", e)
}

func (e *CheckedError) Error() string {
	return fmt.Sprintf("CheckedError : %v", e)
}

func (e *StandardError) initStandardError(message string, err *error) {
	e.message = message
	e.timestamp = time.Now().Unix()
	e.cause = *err
}

// NewStandardError constructor for StandardError
func NewStandardError(message string) *StandardError {
	e := &StandardError{
		code: 000,
	}

	e.initStandardError(message, nil)

	return e
}

// NewStandardErrorWrap constructor for StandardError
func NewStandardErrorWrap(message string, err *error) *StandardError {
	e := &StandardError{
		code: 000,
	}

	e.initStandardError(message, err)

	return e
}

func (e *RuntimeError) initRuntimeError(message string, err *error) {
	e.initStandardError(message, err)

	e.code = 100
}

// NewRuntimeError constructor for RuntimeError
func NewRuntimeError(message string) *RuntimeError {
	e := &RuntimeError{}

	e.initRuntimeError(message, nil)

	return e
}

// NewRuntimeErrorWrap constructor for RuntimeError
func NewRuntimeErrorWrap(message string, err *error) *RuntimeError {
	e := &RuntimeError{}

	e.initRuntimeError(message, err)

	return e
}

func (e *CheckedError) initCheckedError(message string, err *error) {
	e.initStandardError(message, err)

	e.code = 200
}

// NewCheckedError constructor for RuntimeError
func NewCheckedError(message string) *CheckedError {
	e := &CheckedError{}

	e.initCheckedError(message, nil)

	return e
}

// NewCheckedErrorWrap constructor for RuntimeError
func NewCheckedErrorWrap(message string, err *error) *CheckedError {
	e := &CheckedError{}

	e.initCheckedError(message, err)

	return e
}
