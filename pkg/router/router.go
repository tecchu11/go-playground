// Package router implements multiplexer based on http.ServeMux.
package router

import (
	"context"
	"net/http"
)

// Router is wrapper of http.ServeMux.
type Router struct {
	mux                    *http.ServeMux
	notfound               ErrorResponseFunc
	notfoundPattern        string
	methodNotAllowed       ErrorResponseFunc
	methodNotAllowsPattern string
}

type (
	routerOption struct {
		notfound               ErrorResponseFunc
		notfoundPattern        string
		methodNotAllowed       ErrorResponseFunc
		methodNotAllowsPattern string
	}
	// RouterOptionFunc is optional function pattern for Router.
	RouterOptionFunc func(*routerOption)
)

// NotFoundHandler registers given ErrorResponseFunc to router.
func NotFoundHandler(val ErrorResponseFunc) RouterOptionFunc {
	return func(ro *routerOption) {
		ro.notfound = val
	}
}

// NotFoundHandlerPattern is used router pattern when request is dispatched not found handler.
func NotFoundHandlerPattern(val string) RouterOptionFunc {
	return func(ro *routerOption) {
		ro.notfoundPattern = val
	}
}

// MethodNotAllowed registers given ErrorResponseFunc to router.
func MethodNotAllowed(val ErrorResponseFunc) RouterOptionFunc {
	return func(ro *routerOption) {
		ro.methodNotAllowed = val
	}
}

// MethodNotAllowedPattern is used router pattern when request is dispatched method not allowed handler.
func MethodNotAllowedPattern(val string) RouterOptionFunc {
	return func(ro *routerOption) {
		ro.methodNotAllowsPattern = val
	}
}

const (
	defaultNotfoundPattern         = "DefaultNotFoundHandler"
	defaultMethodNotAllowedPattern = "DefaultMethodNotAllowedHandler"
)

// New init Router with given RouterOptionFunc.
func New(options ...RouterOptionFunc) *Router {
	opt := routerOption{
		notfound:               defaultErrorWriter,
		notfoundPattern:        defaultNotfoundPattern,
		methodNotAllowed:       defaultErrorWriter,
		methodNotAllowsPattern: defaultMethodNotAllowedPattern,
	}
	for _, fn := range options {
		fn(&opt)
	}
	return &Router{
		mux:                    http.NewServeMux(),
		notfound:               opt.notfound,
		notfoundPattern:        opt.notfoundPattern,
		methodNotAllowed:       opt.methodNotAllowed,
		methodNotAllowsPattern: opt.methodNotAllowsPattern,
	}
}

var routePatternContextKey struct{}

// ServeHTTP dispatches the request to the handler whose pattern most closely matches the request URL.
//
// The pattern is registered in the request context and can be retrieved from RoutePattern function.
// It will be the registered path matching the request, or in the case of an internally generated redirect, the path to match after the redirect.
//
// This method is based on ServeMux.ServeHTTP. So check ServeMux.ServeHTTP for details.
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	next, pattern := router.mux.Handler(r)
	if pattern == "" {
		// pattern is empty when http.ServeMux determines request is 404 or 405.
		// ServeMux does not provide a way to write custom 404, 405 responses.
		// Therefore, a custom ResponseWriter can be passed to the next handler to return a custom defined response.
		interceptor := newResponseInterceptor(w, r, router.notfound, router.methodNotAllowed)
		ctx := context.WithValue(r.Context(), routePatternContextKey, router.notfoundPattern)
		next.ServeHTTP(interceptor, r.WithContext(ctx))
		return
	}
	next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), routePatternContextKey, pattern)))
}

// Handle registers the handler for the given pattern. If the given pattern conflicts, with one that is already registered, Handle panics.
func (router *Router) Handle(pattern string, handler http.Handler) {
	router.mux.Handle(pattern, handler)
}

// HandleFunc registers the handler function for the given pattern. If the given pattern conflicts, with one that is already registered, HandleFunc panics.
func (router *Router) HandleFunc(pattern string, handler http.HandlerFunc) {
	router.mux.HandleFunc(pattern, handler)
}

// RoutePattern finds pattern registered Router.
func RoutePattern(r *http.Request) string {
	// ignore because pattern is always string.
	pattern := r.Context().Value(routePatternContextKey).(string)
	return pattern
}
