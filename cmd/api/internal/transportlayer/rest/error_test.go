package rest_test

import (
	"go-playground/cmd/api/internal/transportlayer/rest"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErr(t *testing.T) {
	w := httptest.NewRecorder()

	rest.Err(w, "test", http.StatusInternalServerError)

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"message":"test"}`, w.Body.String())
}
