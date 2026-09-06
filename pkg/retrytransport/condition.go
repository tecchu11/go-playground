package retrytransport

import (
	"errors"
	"io"
	"net/http"
)

// ConditionFn determines whether to retry with given http.Response and error.
type ConditionFn func(*http.Response, error) bool

// DefaultCondition is default retry condition determines whether to retry.
// Default is only return when error is io.ErrUnexpectedEOF.
//
// Note: Given error may also be a context cancel or deadline exceeded, so errors should be retried only specific situation.
// Note: Response is nil-able. So you must check whether Response is nil.
var DefaultCondition = ConditionFn(func(_ *http.Response, err error) bool {
	return errors.Is(err, io.ErrUnexpectedEOF)
})
