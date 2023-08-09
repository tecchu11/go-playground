package handler_test

import (
	"encoding/json"
	"go-playground/internal/interactor/rest/handler"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHealthHandler_GetStatus(t *testing.T) {
	test := struct {
		name                string
		inputResponseWriter *httptest.ResponseRecorder
		inputRequest        *http.Request
		expectedCode        int
		expectedBody        map[string]string
	}{
		name:                "test GetStatus returns 200 and expected body",
		inputResponseWriter: httptest.NewRecorder(),
		inputRequest:        httptest.NewRequest("GET", "http://example.com/health", nil),
		expectedCode:        200,
		expectedBody:        map[string]string{"status": "ok"},
	}

	t.Run(test.name, func(t *testing.T) {
		handler.StatusHandler().ServeHTTP(test.inputResponseWriter, test.inputRequest)
		if test.inputResponseWriter.Code != test.expectedCode {
			t.Errorf("unexpected status code(%d) was recieved", test.inputResponseWriter.Code)
		}
		var actualBody map[string]string
		_ = json.Unmarshal(test.inputResponseWriter.Body.Bytes(), &actualBody)
		if !reflect.DeepEqual(actualBody, test.expectedBody) {
			t.Errorf("unexpected body(%v) was recieved", actualBody)
		}
	})
}
