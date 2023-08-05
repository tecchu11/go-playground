package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type newrelicTransactionMiddleWare struct {
	app *newrelic.Application
}

func NewNewrelicTransactionMidleware(app *newrelic.Application) MiddleWare {
	return &newrelicTransactionMiddleWare{app}
}

func (mid *newrelicTransactionMiddleWare) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if mid.app == nil {
			// app is nil means running in local env.
			next.ServeHTTP(w, r)
			return
		}

		txn := mid.app.StartTransaction("")
		defer txn.End()

		txn.SetWebRequestHTTP(r)
		w = txn.SetWebResponse(w)
		ctx := newrelic.NewContext(r.Context(), txn)
		next.ServeHTTP(w, r.WithContext(ctx))
		txn.SetName(txnName(ctx, r.Method))
	})
}


// txnName return transaction name. Transaction name format is "(MethodName) (Pattern)".
// chi defines pattern after handled request and response, so please call this func after serve http.
// Transaction name is "ErrorHandler" when handled by the hadler registered with chi.Mux.NotFound.
func txnName(ctx context.Context, method string) string {
	p := chi.RouteContext(ctx).RoutePattern()
	if p == "" {
		return "ErrorHadnler"
	}
	return method + " " + p
}
