package apperr

import "log/slog"

// Option applies optional info to [Error].
type Option interface {
	do(*Error)
}

type optionFunc func(*Error)

func (fn optionFunc) do(e *Error) {
	fn(e)
}

// WithLevel overrides default log level.
func WithLevel(level slog.Level) Option {
	return optionFunc(func(e *Error) {
		e.logLevel = &level
	})
}

// WithCause wraps given error by [Error].
func WithCause(cause error) Option {
	return optionFunc(func(e *Error) {
		e.cause = cause
	})
}
