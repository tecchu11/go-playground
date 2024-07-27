package handler

import (
	"errors"
	"go-playground/pkg/errorx"
	"go-playground/pkg/problemdetails"
	"log/slog"
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
)

type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (fn ErrorHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	txn := newrelic.FromContext(r.Context())

	err := fn(w, r)
	if err == nil {
		return
	}
	var appErr *errorx.Error
	if !errors.As(err, &appErr) {
		slog.ErrorContext(r.Context(), "unhandled error", slog.String("error", err.Error()))
		txn.NoticeError(err)
		problemdetails.
			New("Unhandled error", http.StatusInternalServerError).
			WithDetail(err.Error()).
			Write(w, r)
		return
	}
	slog.Log(r.Context(), appErr.Level(), "caught error", slog.Any("error", appErr))
	if appErr.Level() == slog.LevelError {
		txn.NoticeError(appErr)
	}
	problemdetails.
		New("Handled error", appErr.HTTPStatus()).
		WithDetail(appErr.Msg()).
		Write(w, r)
}

var _ http.Handler = (ErrorHandlerFunc)(nil)
