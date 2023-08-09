package renderer_test

import (
	"encoding/json"
	"go-playground/pkg/renderer"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestOk(t *testing.T) {
	w := httptest.NewRecorder()
	input := map[string]string{"message": "hello"}
	renderer.Ok(w, input)

	expectedBody := input
	var actualBody map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &actualBody)
	if !reflect.DeepEqual(actualBody, expectedBody) {
		t.Errorf("unexpected body was recived. Actula body is (%v)", actualBody)
	}
	expectedCode := 200
	if actualCode := w.Code; actualCode != expectedCode {
		t.Errorf("http status code %d is unexpected", actualCode)
	}
}
