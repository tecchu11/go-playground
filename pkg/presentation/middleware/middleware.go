package middleware

import (
	"net/http"
)

// MiddleWare is custome interface.
type MiddleWare interface {
	Handle(next http.Handler) http.Handler
}

func Composite(middrewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		for i := range middrewares {
			h = middrewares[len(middrewares)-i-1](h)
		}
		return h
	}
}
