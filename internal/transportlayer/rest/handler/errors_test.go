package handler_test

import (
	"encoding/json"
	"go-playground/internal/transportlayer/rest/handler"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotFoundHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/foos", nil)

	handler.NotFoundHandler(&mockJSON{}).ServeHTTP(w, r)

	expectedCode := 404
	expectedBody := map[string]string{"title": "Resource Not Found", "detail": "/foos resource does not exist"}

	actualCode := w.Code
	var actualBody map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &actualBody)

	assert.NoError(t, err, "json unmarshal should be no err")
	assert.Equal(t, expectedBody, actualBody, "response body should be equal")
	assert.Equal(t, expectedCode, actualCode, "status code should be equal")
}

func TestMethodNotAllowedHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/foos", nil)

	handler.MethodNotAllowedHandler(&mockJSON{}).ServeHTTP(w, r)

	expectedCode := 405
	expectedBody := map[string]string{"title": "Method Not Allowed", "detail": "Http method POST is not allowed for /foos resource"}

	actualCode := w.Code
	var actualBody map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &actualBody)

	assert.NoError(t, err, "json unmarshal should not be err")
	assert.Equal(t, expectedBody, actualBody, "response body should be equal")
	assert.Equal(t, expectedCode, actualCode, "status code should be equal")
}
