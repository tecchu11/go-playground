package errorx

import "log/slog"

func NewErrorForTest(
	msg string,
	cause error,
	at string,
	level slog.Level,
	status int,
) *Error {
	return &Error{
		msg:    msg,
		cause:  cause,
		at:     at,
		level:  level,
		status: status,
	}
}

func (e *Error) Cause() error {
	return e.cause
}

func (e *Error) At() string {
	return e.at
}

func (e *Error) Status() int {
	return e.status
}
