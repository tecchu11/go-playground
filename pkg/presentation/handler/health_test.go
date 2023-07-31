package handler_test

import (
	"encoding/json"
	"go-playground/pkg/presentation/handler"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

func TestHealthHandler_GetStatus(t *testing.T) {
	test := struct {
		name                string
		inputResponseWriter *httptest.ResponseRecorder
		inputRequest        *http.Request
		expectedCode        int
		expectedBody        handler.HealthStatus
	}{
		name:                "test GetStatus returns 200 and expected body",
		inputResponseWriter: httptest.NewRecorder(),
		inputRequest:        httptest.NewRequest("GET", "http://example.com/health", nil),
		expectedCode:        200,
		expectedBody:        handler.HealthStatus{Status: "OK"},
	}

	t.Run(test.name, func(t *testing.T) {
		handler.NewHealthHandler(zap.NewExample()).GetStatus().ServeHTTP(test.inputResponseWriter, test.inputRequest)
		if test.inputResponseWriter.Code != test.expectedCode {
			t.Errorf("unexpected status code(%d) was recieved", test.inputResponseWriter.Code)
		}
		var actualBody handler.HealthStatus
		_ = json.Unmarshal(test.inputResponseWriter.Body.Bytes(), &actualBody)
		if !reflect.DeepEqual(actualBody, test.expectedBody) {
			t.Errorf("unexpected body(%v) was recieved", actualBody)
		}
	})
}