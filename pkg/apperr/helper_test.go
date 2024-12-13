package apperr

import "testing"

// MustAppErr asserts given err by [*Error].
func MustAppErr(t *testing.T, err error) *Error {
	got, ok := (err).(*Error)
	if !ok {
		t.Fatalf("given err is not type of *Error. Actual is %T", err)
	}
	return got
}
