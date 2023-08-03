package middleware

import (
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

		// chi define roter patter after handle request and response, so set transaction name after serve http
		txn.SetName(r.Method + " " + chi.RouteContext(ctx).RoutePattern())
	})
}
