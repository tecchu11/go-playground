package router

type InterceptWriter = interceptWriter

var DefaultErrMarshal = defaultErrMarshal

type PatternContextKey = patternCtxKey

func (router *Router) Middleware() Middleware {
	return router.middleware
}

func (router *Router) Marshal404() ErrMarshalFunc {
	return router.marshal404
}

func (router *Router) Marshal405() ErrMarshalFunc {
	return router.marshal405
}
