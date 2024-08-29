package middleware

import (
	"go-playground/cmd/api/internal/transportlayer/rest"
	"log/slog"
	"net/http"
)

const (
	connectionHeader      = "Connection"
	connectionHeaderValue = "Upgrade"
)

// Recover handle un-recovered panic when handling request.
// If panic have not happened, this middleware nothing to do.
var Recover = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if err == http.ErrAbortHandler {
					// noop
					panic(err)
				}
				slog.ErrorContext(
					r.Context(),
					"Panic was happened. So check detail as soon as possible the reason why happened panic.",
					slog.Any("error", err),
				)
				if r.Header.Get(connectionHeader) != connectionHeaderValue {
					rest.Err(
						w,
						"Unexpected error was happened. Please report this error you have checked.",
						http.StatusInternalServerError,
					)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}
