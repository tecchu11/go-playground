package middleware

import (
	"go-playground/pkg/problemdetails"
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
					panic(err)
				}
				slog.ErrorContext(r.Context(), "Panic was happened. So check detail as soon as possible the reason why happened panic.", slog.Any("error", err))
				if r.Header.Get(connectionHeader) != connectionHeaderValue {
					problemdetails.
						New("Internal Server Error", http.StatusInternalServerError).
						WithDetail("Unexpected error was happened. Please report this error you have checked.").
						Write(w, r)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}
