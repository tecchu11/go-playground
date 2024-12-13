package apperr_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-playground/pkg/apperr"
	"log/slog"
	"net/http"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	type input struct {
		text, msg string
		opts      []apperr.Option
	}
	tests := map[string]struct {
		input input
		want  string
	}{
		"new without options": {
			input: input{text: "this is error info message", msg: "this is message for client"},
			want:  "this is error info message",
		},
		"new with options": {
			input: input{text: "do something", msg: "failed to do something", opts: []apperr.Option{apperr.WithCause(sql.ErrNoRows)}},
			want:  "do something: sql: no rows in result set",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := apperr.New(tc.input.text, tc.input.msg, tc.input.opts...)
			got := apperr.MustAppErr(t, err)

			assert.EqualError(t, got, tc.want)
		})
	}
}

func TestErrorUnwrap(t *testing.T) {
	type input struct {
		text, msg string
		opts      []apperr.Option
	}
	tests := map[string]struct {
		input input
		want  error
	}{
		"unwrapped": {
			input: input{text: "new error on doing something", msg: "failed to do something"},
		},
		"wrapped": {
			input: input{text: "do something", msg: "failed to do something", opts: []apperr.Option{apperr.WithCause(sql.ErrNoRows)}},
			want:  sql.ErrNoRows,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := apperr.New(tc.input.text, tc.input.msg, tc.input.opts...)
			appErr := apperr.MustAppErr(t, err)

			got := appErr.Unwrap()

			if got == nil {
				assert.NoError(t, tc.want)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestErrorLevel(t *testing.T) {
	type input struct {
		text, msg string
		opts      []apperr.Option
	}
	tests := map[string]struct {
		input input
		want  slog.Level
	}{
		"new without options": {
			input: input{text: "text", msg: "msg"},
			want:  slog.LevelError,
		},
		"new with level options": {
			input: input{text: "text", msg: "msg", opts: []apperr.Option{apperr.WithLevel(slog.LevelInfo)}},
			want:  slog.LevelInfo,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := apperr.New(tc.input.text, tc.input.msg, tc.input.opts...)
			appErr := apperr.MustAppErr(t, err)

			got := appErr.Level()

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestErrorCode(t *testing.T) {
	type input struct {
		text, msg string
		opts      []apperr.Option
	}
	tests := map[string]struct {
		input input
		want  apperr.Code
	}{
		"new without options": {
			input: input{text: "text", msg: "msg"},
			want:  apperr.CodeInternal,
		},
		"new with options": {
			input: input{text: "text", msg: "msg", opts: []apperr.Option{apperr.CodeInvalidArgument}},
			want:  apperr.CodeInvalidArgument,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := apperr.New(tc.input.text, tc.input.msg, tc.input.opts...)
			appErr := apperr.MustAppErr(t, err)

			got := appErr.Code()

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestErrorHTTPStatus(t *testing.T) {
	type input struct {
		text, msg string
		opts      []apperr.Option
	}
	tests := map[string]struct {
		input input
		want  int
	}{
		"new without options": {
			input: input{text: "text", msg: "msg"},
			want:  http.StatusInternalServerError,
		},
		"new with options": {
			input: input{text: "text", msg: "msg", opts: []apperr.Option{apperr.CodeInvalidArgument}},
			want:  http.StatusBadRequest,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := apperr.New(tc.input.text, tc.input.msg, tc.input.opts...)
			appErr := apperr.MustAppErr(t, err)

			got := appErr.HTTPStatus()

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestErrorClientMessage(t *testing.T) {
	err := apperr.New("text", "message for client")
	appErr := apperr.MustAppErr(t, err)

	got := appErr.ClientMessage()

	assert.Equal(t, "message for client", got)
}

func TestErrorFormat(t *testing.T) {
	type input struct {
		format    string
		text, msg string
	}
	tests := map[string]struct {
		input input
		want  string
	}{
		"format with %+v": {
			input: input{format: "%+v", text: "text", msg: "msg"},
			want:  `^text\ngo-playground/pkg/apperr_test.TestErrorFormat.func\d+\(.+/pkg/apperr/error_test.go:\d+\)`,
		},
		"format with %v": {
			input: input{format: "%v", text: "text", msg: "msg"},
			want:  `^text\n$`,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := apperr.New(tc.input.text, tc.input.msg)

			got := fmt.Sprintf(tc.input.format, err)

			assert.Regexp(t, regexp.MustCompile(tc.want), got)
		})
	}
}

func TestErrorStackTrace(t *testing.T) {
	e := apperr.New("text", "msg")
	appErr := apperr.MustAppErr(t, e)
	buf := bytes.NewBuffer(nil)
	logger := slog.New(slog.NewJSONHandler(buf, nil))

	logger.Info("test", slog.Any("stacktrace", appErr.StackTrace()))

	type record struct {
		StackTrace string `json:"stacktrace"`
	}
	var got record
	err := json.NewDecoder(buf).Decode(&got)
	require.NoError(t, err)
	assert.Regexp(t, regexp.MustCompile(`^go-playground/pkg/apperr_test.TestErrorStackTrace\(.+/pkg/apperr/error_test.go:\d+\)`), got.StackTrace)
}

func TestIsCode(t *testing.T) {
	type input struct {
		err  error
		code apperr.Code
	}
	tests := map[string]struct {
		input input
		want  bool
	}{
		"err is *apperr.Error(without options)": {
			input: input{err: apperr.New("text", "msg"), code: apperr.CodeInternal},
			want:  true,
		},
		"err is *apperr.Error(with options)": {
			input: input{err: apperr.New("text", "msg", apperr.CodeInvalidArgument), code: apperr.CodeInvalidArgument},
			want:  true,
		},
		"err is not *apperr.Error": {
			input: input{err: sql.ErrNoRows, code: apperr.CodeInternal},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := apperr.IsCode(tc.input.err, tc.input.code)

			assert.Equal(t, tc.want, got)
		})
	}
}
