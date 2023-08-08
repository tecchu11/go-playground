package middleware

import (
	"net/http"
)

// MiddleWare is custome interface.
type MiddleWare interface {
	Handle(next http.Handler) http.Handler
}

