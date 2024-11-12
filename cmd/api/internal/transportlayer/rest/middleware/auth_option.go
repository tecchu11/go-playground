package middleware

type authOption func(*Auth)

func WithSkipRoute(route string) authOption {
	return func(a *Auth) {
		a.skipRoutes[route] = struct{}{}
	}
}
