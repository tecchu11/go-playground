package middleware_test

import (
	"encoding/json"
	"net/http"
)

// mockFailure mock for renderer.Failure
type mockFailure struct{}

func (m *mockFailure) Response(w http.ResponseWriter, _ *http.Request, code int, title string, detail string) {
	res := map[string]string{"title": title, "detail": detail}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(res)
}
