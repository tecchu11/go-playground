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

type recoverMiddleWare struct {
	logger *zap.Logger
}

func NewRecoverMiddleWare(logger *zap.Logger) MiddleWare {
	return &recoverMiddleWare{logger}
}

func (mid *recoverMiddleWare) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				if rec == http.ErrAbortHandler {
					panic(rec)
				}
				mid.logger.Error("Panic was happened. This behavior was unexpected, so check detail as sonn as possible.")
				if r.Header.Get(connectionHeader) != connectionHeaderValue {
					render.InternalServerError(w, "unexpected error was happened, so plese report this error you have checked.", r.URL.Path)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}
