package middleware_test

import (
	"encoding/json"
	"go-playground/internal/transport_layer/rest/middleware"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

func Test_RecoverMiddleWare_Handle(t *testing.T) {
	tests := map[string]struct {
		inputWriter     *httptest.ResponseRecorder
		inputRequest    *http.Request
		expectErr       bool
		expectedCode    int
		expectedBody    map[string]string
		expectedErrBody map[string]string
	}{
		"return 500 and expect body when panic": {
			inputWriter:     httptest.NewRecorder(),
			inputRequest:    httptest.NewRequest("GET", "http://example.com/foos", nil),
			expectErr:       true,
			expectedCode:    500,
			expectedErrBody: map[string]string{"title": "Internal Server Error", "detail": "Unexpected error was happened. Please report this error you have checked."},
		},
		"recover middleware nothing to do": {
			inputWriter:  httptest.NewRecorder(),
			inputRequest: httptest.NewRequest("GET", "http://example.com/foos", nil),
			expectErr:    false,
			expectedCode: 200,
			expectedBody: map[string]string{"hello": "world"},
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			rec := middleware.Recover(zap.NewExample(), &mockFailure{})
			fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if v.expectErr {
					panic("test panic!!")
				}
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(map[string]string{"hello": "world"})
			})
			rec(fn).ServeHTTP(v.inputWriter, v.inputRequest)

			if v.expectErr {
				if actual := v.inputWriter.Code; actual != v.expectedCode {
					t.Errorf("actual status code %d is unexpected. expected is %d", actual, v.expectedCode)
				}
				var actual map[string]string
				_ = json.Unmarshal(v.inputWriter.Body.Bytes(), &actual)
				if !reflect.DeepEqual(actual, v.expectedErrBody) {
					t.Errorf("actual body (%v) is unexpected. expected is (%v)", actual, v.expectedErrBody)
				}
				return
			}

			if actual := v.inputWriter.Code; actual != v.expectedCode {
				t.Errorf("actual status code %d is unexpected. expected is %d", actual, v.expectedCode)
			}
			var actual map[string]string
			_ = json.Unmarshal(v.inputWriter.Body.Bytes(), &actual)
			if !reflect.DeepEqual(actual, v.expectedBody) {
				t.Errorf("actual body (%v) is unexpected. expected is (%v)", actual, v.expectedBody)
			}
		})
	}
}
