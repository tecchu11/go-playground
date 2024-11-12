package ctxhelper

import "context"

type authSubjectKey struct{}

// WithSubject attaches user auth subject into [context.Context].
func WithSubject(ctx context.Context, sub string) context.Context {
	return context.WithValue(ctx, authSubjectKey{}, sub)
}

// Subject retrieves user auth subject from [context.Context].
func Subject(ctx context.Context) (string, bool) {
	sub, ok := ctx.Value(authSubjectKey{}).(string)
	return sub, ok
}
