package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// NewrelicTxn start transaction if newrelic.Application is not nil.
// If newrelic.Application is nil, this middleware nothing to do.
func NewrelicTxn(app *newrelic.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if app != nil {
				next.ServeHTTP(w, r)
				return
			}

			txn := app.StartTransaction("")
			defer txn.End()

			txn.SetWebRequestHTTP(r)
			w = txn.SetWebResponse(w)
			ctx := newrelic.NewContext(r.Context(), txn)
			next.ServeHTTP(w, r.WithContext(ctx))
			txn.SetName(txnName(ctx, r.Method))
		})
		return fn
	}
}

// txnName return transaction name. Transaction name format is "(MethodName) (Pattern)".
// chi defines pattern after handled request and response, so please call this func after serve http.
// Transaction name is "ErrorHandler" when handled by the handler registered with chi.Mux.NotFound.
func txnName(ctx context.Context, method string) string {
	p := chi.RouteContext(ctx).RoutePattern()
	if p == "" {
		return "ErrorHandler"
	}
	return method + " " + p
}

// RequestID is purpose for response to client.
func RequestID(ctx context.Context) string {
	txn := newrelic.FromContext(ctx)
	if txn == nil {
		return ""
	}
	return txn.GetLinkingMetadata().TraceID
}
