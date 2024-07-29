// Package nrhttp provides trace middleware for newrelic.
package nrhttp

import (
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// Middleware traces with newrelic.
func Middleware(app *newrelic.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if app == nil {
				next.ServeHTTP(w, r)
				return
			}
			pattern := r.Pattern
			if pattern == "" {
				pattern = "ErrorHandler" // error handler routes 404 or 405 handler.
			}
			txn := app.StartTransaction(r.Pattern)
			defer txn.End()
			txn.SetWebRequestHTTP(r)
			w = txn.SetWebResponse(w)
			ctx := newrelic.NewContext(r.Context(), txn)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
		return fn
	}
}
