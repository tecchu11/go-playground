package middleware

import (
	"go-playground/pkg/renderer"
	"log/slog"
	"net/http"
)

const (
	connectionHeader      = "Connection"
	connectionHeaderValue = "Upgrade"
)

// Recover handle un-recovered panic when handling request.
// If panic have not happened, this middleware nothing to do.
func Recover(failure renderer.JSON) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					if rec == http.ErrAbortHandler {
						panic(rec)
					}
					slog.ErrorContext(r.Context(), "Panic was happened. So check detail as soon as possible the reason why happened panic.", slog.Any("error", rec))
					if r.Header.Get(connectionHeader) != connectionHeaderValue {
						failure.Failure(w, r, http.StatusInternalServerError, "Internal Server Error", "Unexpected error was happened. Please report this error you have checked.")
					}
				}
			}()
			next.ServeHTTP(w, r)
		})
		return fn
	}
}
