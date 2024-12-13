package apperr

import (
	"fmt"
	"log/slog"
	"net/http"
)

// Code represents application error code.
//
// Code implements [Option] interface.
// So, you can specify code with [New] function.
//
// Code, each of which is mapped to [slog.Level].
// If you want to change the default level, use [WithLevel] option.
type Code uint32

const (
	CodeUnknown Code = iota
	CodeInvalidArgument
	CodeUnAuthn
	CodeUnAuthz
	CodeNotFound
	CodeInternal
)

// String implements [fmt.Stringer].
func (c Code) String() string {
	switch c {
	case CodeUnknown:
		return "unknown"
	case CodeInvalidArgument:
		return "invalidArgument"
	case CodeUnAuthn:
		return "unauthenticated"
	case CodeUnAuthz:
		return "unauthorized"
	case CodeNotFound:
		return "notfound"
	case CodeInternal:
		return "internalServerError"
	default:
		return fmt.Sprintf("Code(%d)", c)
	}
}

func (c Code) status() int {
	switch c {
	case CodeInvalidArgument:
		return http.StatusBadRequest
	case CodeUnAuthn:
		return http.StatusUnauthorized
	case CodeUnAuthz:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	default: // when CodeUnknown, CodeInternal or etc.
		return http.StatusInternalServerError
	}
}

func (c Code) level() slog.Level {
	switch c {
	case CodeInvalidArgument, CodeUnAuthn, CodeUnAuthz, CodeNotFound:
		return slog.LevelWarn
	default:
		return slog.LevelError
	}
}

// do implements [Option] interface.
func (c Code) do(e *Error) {
	e.code = c
}
