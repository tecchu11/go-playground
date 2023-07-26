package middleware

import "net/http"

type headMiddleWare struct {
}

func NewHeadMiddleWare() MiddleWare {
	return &headMiddleWare{}
}

func (mid *headMiddleWare) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
