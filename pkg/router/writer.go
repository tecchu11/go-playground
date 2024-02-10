package router

import (
	"net/http"
)

type ErrorResponseFunc func(http.ResponseWriter, *http.Request, int) (int, error)

var defaultErrorWriter = ErrorResponseFunc(func(w http.ResponseWriter, _ *http.Request, statusCode int) (int, error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	return w.Write(nil)
})

type responseInterceptor struct {
	origin           http.ResponseWriter
	r                *http.Request
	notFound         ErrorResponseFunc
	methodNotAllowed ErrorResponseFunc
	code             int
}

func newResponseInterceptor(origin http.ResponseWriter, r *http.Request, notfound ErrorResponseFunc, methodNotAllowed ErrorResponseFunc) *responseInterceptor {
	if notfound == nil {
		notfound = defaultErrorWriter
	}
	if methodNotAllowed == nil {
		methodNotAllowed = defaultErrorWriter
	}
	return &responseInterceptor{
		origin:           origin,
		r:                r,
		notFound:         notfound,
		methodNotAllowed: methodNotAllowed,
		code:             http.StatusOK,
	}
}

func (w *responseInterceptor) Header() http.Header {
	return w.origin.Header()
}

func (w *responseInterceptor) Write(data []byte) (int, error) {
	switch code := w.code; code {
	case http.StatusNotFound:
		return w.notFound(w.origin, w.r, http.StatusNotFound)
	case http.StatusMethodNotAllowed:
		return w.methodNotAllowed(w.origin, w.r, http.StatusMethodNotAllowed)
	default:
		return w.origin.Write(data)
	}
}

func (w *responseInterceptor) WriteHeader(statusCode int) {
	w.code = statusCode
	if w.code == http.StatusNotFound || w.code == http.StatusMethodNotAllowed {
		// ignore because will write header in ErrorResponseFunc.
		return
	}
	w.origin.WriteHeader(statusCode)
}

var _ http.ResponseWriter = (*responseInterceptor)(nil)
