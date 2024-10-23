package handler

import (
	"errors"
	"go-playground/cmd/api/internal/transportlayer/rest"
	"go-playground/pkg/errorx"
	"log/slog"
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// ErrorHandlerFunc handles request/response safely.
func ErrorHandlerFunc(w http.ResponseWriter, r *http.Request, fn func(http.ResponseWriter, *http.Request) error) {
	txn := newrelic.FromContext(r.Context())

	err := fn(w, r)
	if err == nil {
		return
	}
	var appErr *errorx.Error
	if !errors.As(err, &appErr) {
		slog.ErrorContext(r.Context(), "caught unhandled error", slog.String("error", err.Error()))
		txn.NoticeError(err)
		rest.Err(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Log(r.Context(), appErr.Level(), "caught error", slog.Any("error", appErr))
	if appErr.Level() == slog.LevelError {
		txn.NoticeError(appErr)
	} else {
		txn.NoticeExpectedError(appErr)
	}
	rest.Err(w, appErr.Msg(), appErr.HTTPStatus())
}
