package nrmux

import (
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
)

type WrapResponseWriter = wrapResponseWriter

func (w *WrapResponseWriter) Code() int {
	return w.code
}

func (w *WrapResponseWriter) SetCode(code int) {
	w.code = code
}

func (w *WrapResponseWriter) SetReq(req *http.Request) {
	w.req = req
}

func (w *WrapResponseWriter) SetMarshal404(fn func(*http.Request) ([]byte, error)) {
	w.marshal404 = fn
}

func (w *WrapResponseWriter) SetMarshal405(fn func(*http.Request) ([]byte, error)) {
	w.marshal405 = fn
}

type Option = option

func (o *Option) MarshalJSON404() func(*http.Request) ([]byte, error) {
	return o.marshalJSON404
}

func (o *Option) MarshalJSON405() func(*http.Request) ([]byte, error) {
	return o.marshalJSON405
}

func (mux *NRServeMux) App() *newrelic.Application {
	return mux.app
}

func (mux *NRServeMux) Wrap() func(http.Handler) http.Handler {
	return mux.wrap
}

func (mux *NRServeMux) Unwrap() func(http.Handler) http.Handler {
	return mux.unwrap
}
