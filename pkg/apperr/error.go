package apperr

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

// Error is represents application error.
type Error struct {
	text, msg  string
	code       Code
	cause      error
	logLevel   *slog.Level
	stacktrace stacktrace
}

// New creates new error instance with given info.
func New(text, msg string, opts ...Option) error {
	e := Error{
		text:       text,
		msg:        msg,
		code:       CodeInternal,
		stacktrace: caller(5, 2),
	}
	for _, opt := range opts {
		opt.do(&e)
	}
	return &e
}

// Error is implementation of error interface.
func (e *Error) Error() string {
	if e.cause == nil {
		return e.text
	}
	return fmt.Sprintf("%s: %v", e.text, e.cause)
}

// Unwrap unwraps error.
func (e *Error) Unwrap() error {
	return e.cause
}

// Level is error's severity by [slog.Level].
//
// The default value for Level is determined by [Code].
// LevelError if the [Code] option is not specified.
func (e *Error) Level() slog.Level {
	if e.logLevel != nil {
		return *e.logLevel
	}
	return e.code.level()
}

// Code is error's code.
//
// Default code is [CodeInternal].
func (e *Error) Code() Code {
	return e.code
}

// HTTPStatus is error status code.
//
// Default http status is 500.
func (e *Error) HTTPStatus() int {
	return e.code.status()
}

// ClientMessage is the message for client.
func (e *Error) ClientMessage() string {
	return e.msg
}

// Format implements [fmt.Formatter.Format].
func (e *Error) Format(s fmt.State, verb rune) {
	var str strings.Builder
	str.WriteString(e.Error())
	str.WriteString("\n")
	if verb == 'v' && s.Flag('+') {
		str.WriteString(e.stacktrace.String())
	}
	_, _ = s.Write([]byte(str.String()))
}

// StackTrace returns stack trace.
// This method is useful for logging with slog. For example
//
//	var appErr *apperr.Error
//	ok := errors.As(err, &appErr)
//	if ok {
//		slog.Error("error with stack trace", slog.Any("stacktrace", appErr.StackTrace()))
//	}
func (e *Error) StackTrace() stacktrace {
	return e.stacktrace
}

// IsCode checks given error's code is matched given code.
func IsCode(err error, code Code) bool {
	var appErr *Error
	if ok := errors.As(err, &appErr); ok {
		if appErr.code == code {
			return true
		}
	}
	return false
}

var (
	_ error         = (*Error)(nil)
	_ fmt.Formatter = (*Error)(nil)
)
