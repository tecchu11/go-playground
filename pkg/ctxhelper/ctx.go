package ctxhelper

import (
	"context"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v3"
	"github.com/auth0/go-jwt-middleware/v3/core"
	"github.com/auth0/go-jwt-middleware/v3/validator"
)

// WithSubject attaches user auth subject into [context.Context].
func WithSubject(ctx context.Context, sub string) context.Context {
	return core.SetClaims(ctx, &validator.ValidatedClaims{RegisteredClaims: validator.RegisteredClaims{Subject: sub}})
}

// Subject retrieves user auth subject from [context.Context].
func Subject(ctx context.Context) (sub string, ok bool) {
	claim, err := jwtmiddleware.GetClaims[*validator.ValidatedClaims](ctx)
	if err != nil {
		return
	}
	return claim.RegisteredClaims.Subject, true
}
