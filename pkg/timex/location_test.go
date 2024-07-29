package timex_test

import (
	"go-playground/pkg/timex"
	"testing"
)

func TestJST(t *testing.T) {
	if jst := timex.JST().String(); jst != "Asia/Tokyo" {
		t.Fatal(jst)
	}
}
