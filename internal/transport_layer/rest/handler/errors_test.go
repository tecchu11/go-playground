package handler_test

import (
	"encoding/json"
	"go-playground/internal/transport_layer/rest/handler"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_NotFoundHandler(t *testing.T) {
	tests := map[string]struct {
		testTarget   reflect.Value
		inputWriter  *httptest.ResponseRecorder
		inputRequest *http.Request
		expectedCode int
		expectedBody map[string]string
	}{
		"return 404 and expected body": {
			reflect.ValueOf(handler.NotFoundHandler(&mockFailure{}).ServeHTTP),
			httptest.NewRecorder(),
			httptest.NewRequest("GET", "http://example.com/foos", nil),
			404,
			map[string]string{"title": "Resource Not Found", "detail": "/foos resource does not exist"},
		},
		"return 405 and expected body": {
			reflect.ValueOf(handler.MethodNotAllowedHandler(&mockFailure{}).ServeHTTP),
			httptest.NewRecorder(),
			httptest.NewRequest("GET", "http://example.com/foos", nil),
			405,
			map[string]string{"title": "Method Not Allowed", "detail": "Http method GET is not allowed for /foos resource"},
		},
	}
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			arg1 := reflect.ValueOf(v.inputWriter)
			arg2 := reflect.ValueOf(v.inputRequest)
			v.testTarget.Call([]reflect.Value{arg1, arg2})
			if actual := v.inputWriter.Code; actual != v.expectedCode {
				t.Errorf("actual code %d is difference from expected code", actual)
			}
			var actual map[string]string
			_ = json.Unmarshal(v.inputWriter.Body.Bytes(), &actual)
			if !reflect.DeepEqual(actual, v.expectedBody) {
				t.Errorf("actual response body (%v) is different from expected", actual)
			}
		})
	}
}
