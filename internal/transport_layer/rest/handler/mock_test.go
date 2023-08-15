package handler_test

import (
	"encoding/json"
	"net/http"
)

// mockJSON mock for renderer.Failure
type mockJSON struct{}

func (m *mockJSON) Failure(w http.ResponseWriter, _ *http.Request, code int, title string, detail string) {
	res := map[string]string{"title": title, "detail": detail}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(res)
}

func (m *mockJSON) Success(w http.ResponseWriter, code int, body any) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if code == http.StatusNoContent || code == http.StatusResetContent {
		return
	}
	_ = json.NewEncoder(w).Encode(body)
}
