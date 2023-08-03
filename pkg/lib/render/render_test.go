package render_test

import (
	"encoding/json"
	"go-playground/pkg/lib/render"
	"net/http/httptest"
	"reflect"
	"testing"
)

type testResponse struct {
	Msg string `json:"msg"`
}

const (
	testTitle = "https://github.com/tecchu11/go-playground"
)

func TestOk(t *testing.T) {
	tests := []struct {
		name                string
		inputResponseWriter *httptest.ResponseRecorder
		inputBody           testResponse
		expectedCode        int
		expectedContentType string
		expectedBody        testResponse
	}{
		{
			name:                "test Ok returns 200 and expected body",
			inputResponseWriter: httptest.NewRecorder(),
			inputBody:           testResponse{Msg: "this is test response"},
			expectedCode:        200,
			expectedContentType: "application/json; charset=utf-8",
			expectedBody:        testResponse{Msg: "this is test response"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			render.Ok(test.inputResponseWriter, &test.inputBody)
			if test.inputResponseWriter.Code != test.expectedCode {
				t.Errorf("unexpected status code(%v) was returned", test.inputResponseWriter.Code)
			}
			if test.inputResponseWriter.Header().Get("Content-Type") != test.expectedContentType {
				t.Errorf("unexpected conten type(%s) was recived", test.inputResponseWriter.Header().Get("Content-Type"))
			}
			var actualBody testResponse
			_ = json.Unmarshal(test.inputResponseWriter.Body.Bytes(), &actualBody)
			if !reflect.DeepEqual(actualBody, test.expectedBody) {
				t.Errorf("unmarshled response body was unexpected value (%v)", actualBody)
			}
		})
	}
}

func TestAllStatusFailuer(t *testing.T) {
	tests := []struct {
		name                string
		testFunc            reflect.Value
		inputResponseWriter *httptest.ResponseRecorder
		inputDetail         string
		inputPath           string
		expectedCode        int
		expectedContentType string
		expectedBody        render.ProblemDetail
	}{
		{
			name:                "test Unauthorized returns 401 and expected body",
			testFunc:            reflect.ValueOf(render.Unauthorized),
			inputResponseWriter: httptest.NewRecorder(),
			inputDetail:         "authentication failed",
			inputPath:           "/foos",
			expectedCode:        401,
			expectedContentType: "application/json; charset=utf-8",
			expectedBody: render.ProblemDetail{
				Type:    testTitle,
				Title:   "Unauthorized",
				Detail:  "authentication failed",
				Instant: "/foos",
			},
		},
		{
			name:                "test Forbidden returns 403 and expected body",
			testFunc:            reflect.ValueOf(render.Forbidden),
			inputResponseWriter: httptest.NewRecorder(),
			inputDetail:         "authorization failed",
			inputPath:           "/fuzs",
			expectedCode:        403,
			expectedContentType: "application/json; charset=utf-8",
			expectedBody: render.ProblemDetail{
				Type:    testTitle,
				Title:   "Forbidden",
				Detail:  "authorization failed",
				Instant: "/fuzs",
			},
		},
		{
			name:                "test NotFound returns 404 and expected body",
			testFunc:            reflect.ValueOf(render.NotFound),
			inputResponseWriter: httptest.NewRecorder(),
			inputDetail:         "no resources",
			inputPath:           "/bars",
			expectedCode:        404,
			expectedContentType: "application/json; charset=utf-8",
			expectedBody: render.ProblemDetail{
				Type:    testTitle,
				Title:   "Resource Not Found",
				Detail:  "no resources",
				Instant: "/bars",
			},
		},
		{
			name:                "test MethodNotAllowed returns 405 and expected body",
			testFunc:            reflect.ValueOf(render.MethodNotAllowed),
			inputResponseWriter: httptest.NewRecorder(),
			inputDetail:         "method not allowed",
			inputPath:           "/bars",
			expectedCode:        405,
			expectedContentType: "application/json; charset=utf-8",
			expectedBody: render.ProblemDetail{
				Type:    testTitle,
				Title:   "Method Not Allowed",
				Detail:  "method not allowed",
				Instant: "/bars",
			},
		},
		{
			name:                "test InternalServerError returns 500 and expected body",
			testFunc:            reflect.ValueOf(render.InternalServerError),
			inputResponseWriter: httptest.NewRecorder(),
			inputDetail:         "server error",
			inputPath:           "/bazs",
			expectedCode:        500,
			expectedContentType: "application/json; charset=utf-8",
			expectedBody: render.ProblemDetail{
				Type:    testTitle,
				Title:   "Internal Server Error",
				Detail:  "server error",
				Instant: "/bazs",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			arg1 := reflect.ValueOf(test.inputResponseWriter)
			arg2 := reflect.ValueOf(test.inputDetail)
			arg3 := reflect.ValueOf(test.inputPath)
			test.testFunc.Call([]reflect.Value{arg1, arg2, arg3})

			if test.inputResponseWriter.Code != test.expectedCode {
				t.Errorf("unexpected status code(%v) was returned", test.inputResponseWriter.Code)
			}
			if test.inputResponseWriter.Header().Get("Content-Type") != test.expectedContentType {
				t.Errorf("unexpected conten type(%s) was recived", test.inputResponseWriter.Header().Get("Content-Type"))
			}
			var actualBody render.ProblemDetail
			_ = json.Unmarshal(test.inputResponseWriter.Body.Bytes(), &actualBody)
			if !reflect.DeepEqual(actualBody, test.expectedBody) {
				t.Errorf("unmarshled response body was unexpected value (%v)", actualBody)
			}
		})
	}
}
