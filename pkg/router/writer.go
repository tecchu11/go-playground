package router

import (
	"log/slog"
	"net/http"
)

type interceptWriter struct {
	http.ResponseWriter
	Req        *http.Request
	Code       int
	Marshal404 ErrMarshalFunc
	Marshal405 ErrMarshalFunc
}

var defaultErrMarshal = func(_ *http.Request) ([]byte, error) {
	return nil, nil
}

// Writes only ignores given data when http status are 404 or 405 and then write with marshal404 or marshal405.
func (w *interceptWriter) Write(data []byte) (int, error) {
	w.ResponseWriter.WriteHeader(w.Code)
	switch code := w.Code; code {
	case http.StatusNotFound:
		buf, err := w.Marshal404(w.Req)
		if err != nil {
			slog.ErrorContext(w.Req.Context(), "serveMux router defined missing route and then marshal 404 response body", slog.String("error", err.Error()))
			return w.ResponseWriter.Write(nil)
		}
		return w.ResponseWriter.Write(buf)
	case http.StatusMethodNotAllowed:
		buf, err := w.Marshal405(w.Req)
		if err != nil {
			slog.ErrorContext(w.Req.Context(), "serveMux router defined missing route and then marshal 405 response body", slog.String("error", err.Error()))
			return w.ResponseWriter.Write(nil)
		}
		return w.ResponseWriter.Write(buf)
	default:
		return w.ResponseWriter.Write(data)
	}
}

// WriteHeader stores statusCode into struct filed.
func (w *interceptWriter) WriteHeader(statusCode int) {
	w.Code = statusCode
}

var _ http.ResponseWriter = (*interceptWriter)(nil)
