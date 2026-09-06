package retrytransport

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type RetryTransport struct {
	origin     http.RoundTripper
	maxAttempt int
	condition  ConditionFn
	backoff    BackOffFn
}

func New(origin http.RoundTripper, opts ...OptionFunc) *RetryTransport {
	if origin == nil {
		origin = http.DefaultTransport
	}
	rt := &RetryTransport{
		origin:     origin,
		maxAttempt: 1,
		condition:  DefaultCondition,
		backoff:    DefaultBackOff,
	}
	for _, optFn := range opts {
		optFn(rt)
	}
	return rt
}

var _ http.RoundTripper = (*RetryTransport)(nil)

func (rt *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var (
		bodyBuf    []byte
		err        error
		res        *http.Response
		retryCount int
	)
	if req.Body != nil {
		bodyBuf, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("io.ReadAll with non nil request body: %w", err)
		}
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBuf))
	}
	res, err = rt.origin.RoundTrip(req)
loop:
	for {
		if !rt.condition(res, err) || retryCount >= rt.maxAttempt {
			break loop
		}
		bf := rt.backoff(retryCount)
		slog.DebugContext(req.Context(), "Retrying request",
			slog.String("targetHost", req.Host),
			slog.Int("retryCount", retryCount),
			slog.Duration("backOff", bf),
		)
		select {
		case <-time.After(bf):
			drainBody(res)
			if req.Body != nil {
				req.Body = io.NopCloser(bytes.NewBuffer(bodyBuf))
			}
			res, err = rt.origin.RoundTrip(req)
			retryCount++
		case <-req.Context().Done():
			break loop
		}
	}
	return res, err
}

// drainBody discards and closes response body to reuse connection.
func drainBody(res *http.Response) {
	if res != nil && res.Body != nil {
		_, _ = io.Copy(io.Discard, res.Body)
		res.Body.Close()
	}
}
