package handler_test

import (
	"context"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/cmd/api/internal/transportlayer/rest/handler/v2"
	"go-playground/pkg/apperr"
	"go-playground/pkg/ctxhelper"
	"go-playground/pkg/testhelper"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestUserHandler_GetMe(t *testing.T) {
	type input struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	type setup func(*testing.T) (input, *handler.UserHandler)
	type want struct {
		status int
		body   string
	}
	tests := map[string]struct {
		setup setup
		want  want
	}{
		"success to find me": {
			setup: func(t *testing.T) (input, *handler.UserHandler) {
				sub := "01951972-20ce-7002-b4d7-48ca29ffef2b"
				ctx := ctxhelper.WithSubject(context.Background(), sub)
				mck := new(MockUserInteractor)
				mck.On("FindBySub", mock.Anything, sub).Return(entity.User{
					ID:            testhelper.UUIDFromString(t, "01951975-21cb-7991-b0a4-a854c842e258"),
					Sub:           sub,
					FamilyName:    "family",
					GivenName:     "given",
					Email:         "Alvis.Cummerata@example.com",
					EmailVerified: true,
					CreatedAt:     time.Date(2025, 02, 18, 14, 0, 0, 0, time.UTC),
					UpdatedAt:     time.Date(2025, 02, 18, 14, 0, 0, 0, time.UTC),
				}, nil)
				return input{
					r: httptest.NewRequestWithContext(ctx, http.MethodGet, "/users/me", nil),
					w: httptest.NewRecorder(),
				}, &handler.UserHandler{UserInteractor: mck}
			},
			want: want{
				status: http.StatusOK,
				body: `
				{
					"id":"01951975-21cb-7991-b0a4-a854c842e258",
					"sub":"01951972-20ce-7002-b4d7-48ca29ffef2b",
					"familyName":"family",
					"givenName":"given",
					"email":"Alvis.Cummerata@example.com",
					"emailVerified":true,
					"createdAt":"2025-02-18T14:00:00Z",
					"updatedAt":"2025-02-18T14:00:00Z"
				}
				`,
			},
		},
		"failed to find me": {
			setup: func(t *testing.T) (input, *handler.UserHandler) {
				sub := "01951972-20ce-7002-b4d7-48ca29ffef2b"
				ctx := ctxhelper.WithSubject(context.Background(), sub)
				mck := new(MockUserInteractor)
				mck.On("FindBySub", mock.Anything, sub).Return(entity.User{}, apperr.New("fail", "error message", apperr.CodeNotFound))
				return input{
					r: httptest.NewRequestWithContext(ctx, http.MethodGet, "/users/me", nil),
					w: httptest.NewRecorder(),
				}, &handler.UserHandler{UserInteractor: mck}
			},
			want: want{
				status: http.StatusNotFound,
				body:   `{"message":"error message"}`,
			},
		},
		"missing subject": {
			setup: func(t *testing.T) (input, *handler.UserHandler) {
				return input{
					r: httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/users/me", nil),
					w: httptest.NewRecorder(),
				}, &handler.UserHandler{}
			},
			want: want{
				status: http.StatusForbidden,
				body:   `{"message":"authorization failure"}`,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			input, h := tc.setup(t)

			h.GetMe(input.w, input.r)

			assert.JSONEq(t, tc.want.body, input.w.Body.String())
			assert.Equal(t, tc.want.status, input.w.Code)
		})
	}
}
