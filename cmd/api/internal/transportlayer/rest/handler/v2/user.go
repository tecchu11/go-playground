package handler

import (
	"encoding/json"
	"go-playground/cmd/api/internal/transportlayer/rest"
	"go-playground/cmd/api/internal/transportlayer/rest/oapi"
	"go-playground/pkg/ctxhelper"
	"go-playground/pkg/errorx"
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
			rest.Err(w, "missing authenticated user info", http.StatusBadRequest)
			return nil
		}
		var body oapi.RequestUser
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			return errorx.NewWarn("failed to unmarshal request body of PostUser", errorx.WithCause(err), errorx.WithStatus(http.StatusBadRequest))
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
