package errorx

type OptionFunc func(*Error)

// WithCause configures caused error.
func WithCause(err error) OptionFunc {
	return func(e *Error) {
		e.cause = err
	}
}

// WithStatus configures http status.
func WithStatus(status int) OptionFunc {
	return func(e *Error) {
		e.status = status
	}
}
