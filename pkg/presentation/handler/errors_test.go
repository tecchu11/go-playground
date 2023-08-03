package handler_test

import (
	"encoding/json"
	"go-playground/pkg/lib/render"
	"go-playground/pkg/presentation/handler"
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
		expectedBody render.ProblemDetail
	}{
		"return 404 and expected body": {
			reflect.ValueOf(handler.NotFoundHandler().ServeHTTP),
			httptest.NewRecorder(),
			httptest.NewRequest("GET", "http://example.com/foos", nil),
			404,
			render.ProblemDetail{
				Type:    "https://github.com/tecchu11/go-playground",
				Title:   "Resource Not Found",
				Detail:  "A request for a resource that does not exist.",
				Instant: "/foos",
			},
		},
		"return 405 and expected body": {
			reflect.ValueOf(handler.MethodNotAllowedHandler().ServeHTTP),
			httptest.NewRecorder(),
			httptest.NewRequest("GET", "http://example.com/foos", nil),
			405,
			render.ProblemDetail{
				Type:    "https://github.com/tecchu11/go-playground",
				Title:   "Method Not Allowed",
				Detail:  "Http method GET is not allowed for this resource.",
				Instant: "/foos",
			},
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
			var actual render.ProblemDetail
			_ = json.Unmarshal(v.inputWriter.Body.Bytes(), &actual)
			if !reflect.DeepEqual(actual, v.expectedBody) {
				t.Errorf("actual response body (%v) is different from expected", actual)
			}
		})
	}
}
