package retrytransport

import (
	"math"
	"time"
)

// BackOffFn determines backoff periods with given retryCount.
// Default implementation is DefaultBackOff.
type BackOffFn func(retryCount int) time.Duration

// DefaultBackOff is backoff policy for retry.
// BackOff duration is (2^retryCount) * 50 msec. Max backoff duration is 200msec.
var DefaultBackOff = BackOffFn(func(retryCount int) time.Duration {
	var backoff float64
	if retryCount >= 3 { // to prevent non-need backoff calculation.
		backoff = 200
	} else {
		backoff = math.Pow(2, float64(retryCount)) * 50
	}
	return time.Duration(backoff) * time.Millisecond
})
