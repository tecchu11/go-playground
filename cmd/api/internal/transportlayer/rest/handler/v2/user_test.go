package handler_test

import (
	"context"
	"go-playground/cmd/api/internal/transportlayer/rest/handler/v2"
	"go-playground/pkg/apperr"
	"go-playground/pkg/ctxhelper"
	"go-playground/pkg/testhelper"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_PostCreate(t *testing.T) {
	type input struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	type setup func(*testing.T) *handler.UserHandler
	type want struct {
		status int
		body   string
	}
	tests := map[string]struct {
		input input
		setup setup
		want  want
	}{
		"success": {
			input: input{
				w: httptest.NewRecorder(),
				r: httptest.NewRequestWithContext(ctxhelper.WithSubject(context.Background(), "sub1"), http.MethodPost, "/user", strings.NewReader(`
		{
		  "givenName": "Dibbert",
		  "familyName": "Kozey",
		  "email": "Jonathan74@example.com",
		  "emailVerified": true
		}
						`)),
			},
			setup: func(t *testing.T) *handler.UserHandler {
				mck := new(MockUserInteractor)
				mck.
					On("CreateUser", ctxhelper.WithSubject(context.Background(), "sub1"), "sub1", "Dibbert", "Kozey", "Jonathan74@example.com", true).
					Return(testhelper.UUIDFromString(t, "0193196b-28c4-7337-a891-e728860339cd"), nil)
				h := handler.UserHandler{UserInteractor: mck}
				return &h
			},
			want: want{
				status: http.StatusOK,
				body:   `{"id":"0193196b-28c4-7337-a891-e728860339cd"}`,
			},
		},
		"failure missing subject from context": {
			input: input{
				w: httptest.NewRecorder(),
				r: httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/user", strings.NewReader("")),
			},
			setup: func(t *testing.T) *handler.UserHandler { return new(handler.UserHandler) },
			want: want{
				status: http.StatusForbidden,
				body:   `{"message":"authorization failure"}`,
			},
		},
		"failure unmarshal request body": {
			input: input{
				w: httptest.NewRecorder(),
				r: httptest.NewRequestWithContext(ctxhelper.WithSubject(context.Background(), "sub1"), http.MethodPost, "/user", strings.NewReader("")),
			},
			setup: func(t *testing.T) *handler.UserHandler { return new(handler.UserHandler) },
			want: want{
				status: http.StatusBadRequest,
				body:   `{"message":"invalid request"}`,
			},
		},
		"failure interactor returns error": {
			input: input{
				w: httptest.NewRecorder(),
				r: httptest.NewRequestWithContext(ctxhelper.WithSubject(context.Background(), "sub1"), http.MethodPost, "/user", strings.NewReader(`
		{
		  "givenName": "Dibbert",
		  "familyName": "Kozey",
		  "emailVerified": true
		}
						`)),
			},
			setup: func(t *testing.T) *handler.UserHandler {
				mck := new(MockUserInteractor)
				mck.
					On("CreateUser", ctxhelper.WithSubject(context.Background(), "sub1"), "sub1", "Dibbert", "Kozey", "", true).
					Return(uuid.Nil, apperr.New("validation error", "email is required", apperr.CodeInvalidArgument))
				h := handler.UserHandler{UserInteractor: mck}
				return &h
			},
			want: want{
				status: http.StatusBadRequest,
				body:   `{"message":"email is required"}`,
			},
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			h := v.setup(t)

			h.PostUser(v.input.w, v.input.r)

			assert.Equal(t, v.want.status, v.input.w.Code)
			assert.JSONEq(t, v.want.body, v.input.w.Body.String())
		})
	}
	_ = tests
}
