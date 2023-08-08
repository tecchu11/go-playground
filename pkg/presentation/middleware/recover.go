package middleware

import (
	"go-playground/pkg/lib/render"
	"net/http"

	"go.uber.org/zap"
)

const (
	connectionHeader      = "Connection"
	connectionHeaderValue = "Upgrade"
)

// Recover handle unrecovered panic wher handle request.
func Recover(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					if rec == http.ErrAbortHandler {
						panic(rec)
					}
					logger.Error("Panic was happend. So check detail as soon as possible the reason why happened pacnic.")
					if r.Header.Get(connectionHeader) != connectionHeaderValue {
						render.InternalServerError(w, "Unexpected error was happened. Plese report this error you have checked.", r.URL.Path)
					}
				}
			}()
			next.ServeHTTP(w, r)
		})
		return fn
	}
}
