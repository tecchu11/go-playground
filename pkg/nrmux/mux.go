package nrmux

import (
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// NRServeMux is router based on http.ServeMux.
type NRServeMux struct {
	*http.ServeMux

	app    *newrelic.Application
	wrap   func(http.Handler) http.Handler
	unwrap func(http.Handler) http.Handler
}

type (
	option struct {
		marshalJSON404 func(*http.Request) ([]byte, error)
		marshalJSON405 func(*http.Request) ([]byte, error)
	}
	OptionFunc func(*option)
)

// WithMarshalJSON404 configures marshaler for NotFound response body.
func WithMarshalJSON404(fn func(*http.Request) ([]byte, error)) OptionFunc {
	return func(o *option) {
		o.marshalJSON404 = fn
	}
}

// WithMarshalJSON405 configures marshaler for MethodNotAllowed response body.
func WithMarshalJSON405(fn func(*http.Request) ([]byte, error)) OptionFunc {
	return func(o *option) {
		o.marshalJSON405 = fn
	}
}

// New initializes NRServeMux.
// By default, NRServeMux writes nil response body on 404 or 405.
func New(app *newrelic.Application, opts ...OptionFunc) *NRServeMux {
	o := option{
		marshalJSON404: func(r *http.Request) ([]byte, error) { return nil, nil },
		marshalJSON405: func(r *http.Request) ([]byte, error) { return nil, nil },
	}
	for _, opt := range opts {
		opt(&o)
	}
	return &NRServeMux{
		ServeMux: http.NewServeMux(),
		app:      app,
		wrap:     wrapWRMiddleware(o.marshalJSON404, o.marshalJSON405),
		unwrap:   unwrapWRMiddleware,
	}
}

// Handle registers handler which is instrumented by NewRelic.
func (mux *NRServeMux) Handle(pattern string, handler http.Handler) {
	p, hn := newrelic.WrapHandle(mux.app, pattern, handler)
	mux.ServeMux.Handle(p, mux.unwrap(hn))
}

// HandleFunc registers handler which is instrumented by NewRelic.
func (mux *NRServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mux.Handle(pattern, http.HandlerFunc(handler))
}

// ServeHTTP dispatches the request to the handler.
func (mux *NRServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.wrap(mux.ServeMux).ServeHTTP(w, r)
}
