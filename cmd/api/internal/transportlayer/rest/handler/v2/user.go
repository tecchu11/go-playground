package handler

import (
	"encoding/json"
	"go-playground/cmd/api/internal/transportlayer/rest/oapi"
	"go-playground/pkg/apperr"
	"go-playground/pkg/ctxhelper"
	"log/slog"
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
)

type UserHandler struct {
	UserInteractor UserInteractor
}

// PostUser create new user for [POST /users]
func (u *UserHandler) PostUser(w http.ResponseWriter, r *http.Request) {
	defer newrelic.FromContext(r.Context()).StartSegment("handler/UserHandler/PostUser").End()

	ErrorHandlerFunc(w, r, func(w http.ResponseWriter, r *http.Request) error {
		sub, ok := ctxhelper.Subject(r.Context())
		if !ok {
			return apperr.New("subject is missing but this is unexpected", "authorization failure", apperr.CodeUnAuthz, apperr.WithLevel(slog.LevelError))
		}
		var body oapi.RequestUser
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			return apperr.New("unmarshal PostUser request body", "invalid request", apperr.WithCause(err), apperr.CodeInvalidArgument)
		}
		uid, err := u.UserInteractor.CreateUser(
			r.Context(),
			sub,
			body.GivenName,
			body.FamilyName,
			string(body.Email),
			body.EmailVerified,
		)
		if err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(oapi.ResponseUserID{
			ID: uid,
		})
	})
}
