package middleware

import (
	"context"
	"go-playground/pkg/router"
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// NewrelicTxn start transaction if newrelic.Application is not nil.
// If newrelic.Application is nil, this middleware nothing to do.
func NewrelicTxn(app *newrelic.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if app == nil {
				next.ServeHTTP(w, r)
				return
			}
			pattern := router.RoutePattern(r)
			txn := app.StartTransaction(pattern)
			defer txn.End()
			txn.SetWebRequestHTTP(r)
			w = txn.SetWebResponse(w)
			ctx := newrelic.NewContext(r.Context(), txn)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
		return fn
	}
}

// RequestID is purpose for response to client.
func RequestID(ctx context.Context) string {
	txn := newrelic.FromContext(ctx)
	if txn == nil {
		return ""
	}
	return txn.GetLinkingMetadata().TraceID
}
