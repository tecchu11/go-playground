package router

import (
	"context"
	"log/slog"
	"net/http"
)

// ErrMarshalFunc marshals 404 or 405 response body with given http.Request.
type ErrMarshalFunc func(*http.Request) ([]byte, error)

type Middleware func(next http.Handler) http.Handler

// Router is a router based on http.ServeMux.
// It uses the routing feature of http.ServeMux,
// but has been extended to return a user-defined response body when the HTTP Status Code is 404 or 405.
type Router struct {
	*http.ServeMux
	middleware Middleware
	marshal404 ErrMarshalFunc
	marshal405 ErrMarshalFunc
}

// Option configures Router.
type Option func(*Router)

// WithMiddleware configures middleware wraps any http.Handler
func WithMiddleware(mid Middleware) Option {
	return func(r *Router) {
		r.middleware = mid
	}
}

// With404Func configures ErrMarshalFunc for 404.
func With404Func(fn ErrMarshalFunc) Option {
	return func(r *Router) {
		r.marshal404 = fn
	}
}

// With405Func configures ErrMarshalFunc for 405.
func With405Func(fn ErrMarshalFunc) Option {
	return func(r *Router) {
		r.marshal405 = fn
	}
}

// New initializes Router with given options.
// By default, 404 and 405 Response Bodies are written as nil.
func New(opts ...Option) *Router {
	router := Router{
		ServeMux:   http.NewServeMux(),
		marshal404: defaultErrMarshal,
		marshal405: defaultErrMarshal,
	}
	for _, opt := range opts {
		opt(&router)
	}
	return &router
}

type patternCtxKey struct{}

// ServeHTTP dispatches the request to the handler whose pattern most closely matches the request URL.
//
// The logic to determines the handler uses Handler method provided by http.ServeMux.
// This also tries to conform to the behavior of ServeMux.ServeHTTP as much as possible.
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	next, pattern := router.ServeMux.Handler(r)
	if pattern == "" {
		iw := &interceptWriter{
			ResponseWriter: w,
			Req:            r,
			Marshal404:     router.marshal404,
			Marshal405:     router.marshal405,
		}
		ctx := context.WithValue(r.Context(), patternCtxKey{}, "MissingRoutePattern")
		if router.middleware != nil {
			router.middleware(next).ServeHTTP(iw, r.WithContext(ctx))
			return
		}
		next.ServeHTTP(iw, r.WithContext(ctx))
		return
	}
	ctx := context.WithValue(r.Context(), patternCtxKey{}, pattern)
	if router.middleware != nil {
		router.middleware(router.ServeMux).ServeHTTP(w, r.WithContext(ctx))
		return
	}
	router.ServeMux.ServeHTTP(w, r.WithContext(ctx))
}

// Pattern finds handler pattern from given http.Request.
// This function is intended to be used, for example, when implementing OpenTelemetry Middleware.
func Pattern(r *http.Request) string {
	v := r.Context().Value(patternCtxKey{})
	if v == nil {
		slog.Error("missing routing pattern from request context")
		return "MissingRoutingPattern"
	}
	pattern, ok := v.(string)
	if !ok {
		slog.Error("missing routing pattern string from request context")
		return "MissingRoutingPatternString"
	}
	return pattern
}
