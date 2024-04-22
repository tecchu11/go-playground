package nrmux

import (
	"log/slog"
	"net/http"
)

// wrapResponseWriter wraps net/http generated and passed http.ResponseWriter to capture status code.
type wrapResponseWriter struct {
	http.ResponseWriter
	code       int
	req        *http.Request
	marshal404 func(*http.Request) ([]byte, error)
	marshal405 func(*http.Request) ([]byte, error)
}

// Unwrap unwraps w and then returned original ResponseWriter.
func (w *wrapResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

// WriteHeader captures status code and then skips to writer header by given status code.
// Status code will be written in Write method.
func (w *wrapResponseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
}

// Write writes header and then writes body with given marshaler.
func (w *wrapResponseWriter) Write(data []byte) (int, error) {
	w.ResponseWriter.Header().Set("Content-Type", "application/json")
	if w.code != 404 && w.code != 405 {
		slog.WarnContext(w.req.Context(), "Captured status code is not 404 and 405. But this is unexpected", slog.Int("capturedCode", w.code))
		w.code = http.StatusOK
	}
	w.ResponseWriter.WriteHeader(w.code)
	switch w.code {
	case http.StatusNotFound:
		buf, err := w.marshal404(w.req)
		if err != nil {
			slog.ErrorContext(w.req.Context(), "writer nil because failed to marshal 404 body", slog.String("error", err.Error()))
			return w.ResponseWriter.Write(nil)
		}
		return w.ResponseWriter.Write(buf)
	case http.StatusMethodNotAllowed:
		buf, err := w.marshal405(w.req)
		if err != nil {
			slog.ErrorContext(w.req.Context(), "writer nil because failed to marshal 405 body", slog.String("error", err.Error()))
			return w.ResponseWriter.Write(nil)
		}
		return w.ResponseWriter.Write(buf)
	default:
		return w.ResponseWriter.Write(data)
	}
}

// wrapWRMiddleware wraps ResponseWriter by wrapResponseWriter.
func wrapWRMiddleware(
	marshal404 func(*http.Request) ([]byte, error),
	marshal405 func(*http.Request) ([]byte, error),
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wwr := wrapResponseWriter{
				ResponseWriter: w,
				req:            r,
				marshal404:     marshal404,
				marshal405:     marshal405,
			}
			next.ServeHTTP(&wwr, r)
		})
	}
}

type unwrapWriter interface {
	Unwrap() http.ResponseWriter
}

// unwrapWRMiddleware unwraps ResponseWriter.
var unwrapWRMiddleware = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for {
			switch t := w.(type) {
			case *wrapResponseWriter:
				original := t.Unwrap()
				next.ServeHTTP(original, r)
				return
			case unwrapWriter:
				w = t.Unwrap()
			default:
				slog.ErrorContext(r.Context(), "ResponseWriter expects to be wrapped, but the type assertion by *wrapResponseWriter failed.")
				next.ServeHTTP(w, r)
				return
			}
		}
	})
}
