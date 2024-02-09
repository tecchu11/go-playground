package router

import (
	"context"
	"net/http"
)

// Router is wrapper of http.ServeMux.
type Router struct {
	mux                    *http.ServeMux
	notfound               http.Handler
	notfoundPattern        string
	methodNotAllowed       http.Handler
	methodNotAllowsPattern string
}

type (
	routerOption struct {
		notfound               http.Handler
		notfoundPattern        string
		methodNotAllowed       http.Handler
		methodNotAllowsPattern string
	}
	// RouterOptionFunc is optional function pattern for Router.
	RouterOptionFunc func(*routerOption)
)

// NotFoundHandler registers given to router.
func NotFoundHandler(val http.Handler) RouterOptionFunc {
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

// MethodNotAllowed registers given handler to router.
func MethodNotAllowed(val http.Handler) RouterOptionFunc {
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

var (
	defaultNotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write(nil)
	})
	defaultMethodNotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(nil)
	})
)

const (
	defaultNotfoundPattern         = "DefaultNotFoundHandler"
	defaultMethodNotAllowedPattern = "DefaultMethodNotAllowedHandler"
)

// New init Router with given RouterOptionFunc.
func New(options ...RouterOptionFunc) *Router {
	opt := routerOption{
		notfound:               defaultNotFoundHandler,
		notfoundPattern:        defaultNotfoundPattern,
		methodNotAllowed:       defaultMethodNotFoundHandler,
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
		// TODO find solution to determine whether next handle is http.NotFoundHandler
		router.notfound.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), routePatternContextKey, router.notfoundPattern)))
		return
		// router.methodNotAllowed.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), routePatternContextKey, router.methodNotAllowsPattern)))
		// return
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
