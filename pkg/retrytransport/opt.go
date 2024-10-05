package retrytransport

type OptionFunc func(*RetryTransport)

func WithMaxAttempt(num int) OptionFunc {
	return func(rt *RetryTransport) {
		rt.maxAttempt = num
	}
}

func WithCondition(cond ConditionFn) OptionFunc {
	return func(rt *RetryTransport) {
		rt.condition = cond
	}
}

func WithBackoff(backoff BackOffFn) OptionFunc {
	return func(rt *RetryTransport) {
		rt.backoff = backoff
	}
}
