package middleware_test

import (
	"encoding/json"
	"go-playground/internal/interactor/rest/middleware"
	"go-playground/pkg/render"
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
		expectedErrBody render.ProblemDetail
	}{
		"return 500 and expect body when panic": {
			inputWriter:  httptest.NewRecorder(),
			inputRequest: httptest.NewRequest("GET", "http://example.com/foos", nil),
			expectErr:    true,
			expectedCode: 500,
			expectedErrBody: render.ProblemDetail{
				Type:    "https://github.com/tecchu11/go-playground",
				Title:   "Internal Server Error",
				Detail:  "Unexpected error was happened. Plese report this error you have checked.",
				Instant: "/foos",
			},
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
			rec := middleware.Recover(zap.NewExample())
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
				var actual render.ProblemDetail
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
