package errorx

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
)

// Error is error implementation.
// This holds information below.
//   - Caused error.
//   - Where the Error occurred.
//   - Severity of errors by slog.Level
//   - Http status(Default 500).
type Error struct {
	msg    string
	cause  error
	at     string
	level  slog.Level
	status int
}

func callerAt() string {
	pc := make([]uintptr, 10)
	n := runtime.Callers(3, pc)
	if n == 0 {
		return ""
	}
	pc = pc[:n]
	var stack string
	frames := runtime.CallersFrames(pc)
	for {
		frame, more := frames.Next()
		stack += fmt.Sprintf("%s(%s:%d)\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
	return stack
}

func newError(msg, at string, level slog.Level, opts ...OptionFunc) *Error {
	err := Error{
		msg:    msg,
		at:     at,
		level:  level,
		status: http.StatusInternalServerError,
	}
	for _, fn := range opts {
		fn(&err)
	}
	return &err
}

// NewInfo inits info level error.
func NewInfo(msg string, opts ...OptionFunc) *Error {
	return newError(msg, callerAt(), slog.LevelInfo, opts...)
}

// NewWarn inits warn level error.
func NewWarn(msg string, opts ...OptionFunc) *Error {
	return newError(msg, callerAt(), slog.LevelWarn, opts...)
}

// NewError inits error level error.
func NewError(msg string, opts ...OptionFunc) *Error {
	return newError(msg, callerAt(), slog.LevelError, opts...)
}

// Error is error.Error implementation.
func (e *Error) Error() string {
	if e.cause == nil {
		return e.msg
	}
	return fmt.Sprintf("%s [cause]%s [at]%s", e.msg, e.cause.Error(), e.at)
}

// LogValue is slog.LogValuer.LogValue implementation.
func (e *Error) LogValue() slog.Value {
	var errAttr, atAttr slog.Attr
	if e.cause != nil {
		errAttr = slog.String("cause", e.cause.Error())
	}
	if e.at != "" {
		atAttr = slog.String("at", e.at)
	}
	return slog.GroupValue(slog.String("msg", e.msg), errAttr, atAttr)
}

// HTTPStatus represent http status.
func (e *Error) HTTPStatus() int {
	return e.status
}

// Msg reports error message.
func (e *Error) Msg() string {
	return e.msg
}

// Level reports error severity by slog.Level.
func (e *Error) Level() slog.Level {
	return e.level
}

var (
	_ error          = (*Error)(nil)
	_ slog.LogValuer = (*Error)(nil)
)
