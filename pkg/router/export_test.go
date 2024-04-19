package router

type InterceptWriter = interceptWriter

var DefaultErrMarshal = defaultErrMarshal

type PatternContextKey = patternCtxKey

func (r *Router) Middleware() Middleware {
	return r.middleware
}

func (r *Router) Marshal404() ErrMarshalFunc {
	return r.marshal404
}

func (r *Router) Marshal405() ErrMarshalFunc {
	return r.marshal405
}
