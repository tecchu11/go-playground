package handler_test

import (
	"encoding/json"
	"go-playground/internal/transport_layer/rest/handler"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/statuses", nil)

	handler.StatusHandler(&mockJSON{}).ServeHTTP(w, r)

	expectedCode := 200
	expectedBody := handler.Status

	actualCode := w.Code
	var actualBody map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &actualBody)

	assert.NoError(t, err, "json unmarshal should not be err")
	assert.Equal(t, actualBody, expectedBody, "response body should be equal")
	assert.Equal(t, actualCode, expectedCode, "status code should be equal")
}
