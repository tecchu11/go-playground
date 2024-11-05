package errorx_test

import (
	"go-playground/pkg/errorx"
	"io"
	"log/slog"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := map[string]struct {
		testTargetFunc func(string, ...errorx.OptionFunc) *errorx.Error
		inMsg          string
		inOptions      []errorx.OptionFunc
		expected       *errorx.Error
	}{
		"[NewInfo]only message": {
			testTargetFunc: errorx.NewInfo,
			inMsg:          "info error",
			expected:       errorx.NewErrorForTest("info error", nil, "", slog.LevelInfo, 500),
		},
		"[NewInfo]message with options": {
			testTargetFunc: errorx.NewInfo,
			inMsg:          "info error",
			inOptions:      []errorx.OptionFunc{errorx.WithCause(io.EOF), errorx.WithStatus(400)},
			expected:       errorx.NewErrorForTest("info error", io.EOF, "", slog.LevelInfo, 400),
		},
		"[NewWarn]only message": {
			testTargetFunc: errorx.NewWarn,
			inMsg:          "warn error",
			expected:       errorx.NewErrorForTest("warn error", nil, "", slog.LevelWarn, 500),
		},
		"[NewWarn]message with options": {
			testTargetFunc: errorx.NewWarn,
			inMsg:          "warn error",
			inOptions:      []errorx.OptionFunc{errorx.WithCause(io.EOF), errorx.WithStatus(400)},
			expected:       errorx.NewErrorForTest("warn error", io.EOF, "", slog.LevelWarn, 400),
		},
		"[NewError]only message": {
			testTargetFunc: errorx.NewError,
			inMsg:          "error error",
			expected:       errorx.NewErrorForTest("error error", nil, "", slog.LevelError, 500),
		},
		"[NewError]message with options": {
			testTargetFunc: errorx.NewError,
			inMsg:          "error error",
			inOptions:      []errorx.OptionFunc{errorx.WithCause(io.EOF), errorx.WithStatus(400)},
			expected:       errorx.NewErrorForTest("error error", io.EOF, "", slog.LevelError, 400),
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			actual := v.testTargetFunc(v.inMsg, v.inOptions...)
			assert.Equal(t, v.expected.Msg(), actual.Msg())
			assert.Equal(t, v.expected.Cause(), actual.Cause())
			assert.NotEmpty(t, actual.At())
			assert.Equal(t, v.expected.Level(), actual.Level())
			assert.Equal(t, v.expected.Status(), actual.Status())
		})
	}
}

func TestErrorError(t *testing.T) {
	tests := map[string]struct {
		inMsg    string
		inCause  error
		expected string
	}{
		"only message": {inMsg: "error!", expected: "error!"},
		"with cause":   {inMsg: "error!", inCause: io.EOF, expected: `error\! \[cause\]EOF \[at\].*`},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			err := errorx.NewError(v.inMsg, errorx.WithCause(v.inCause))
			actual := err.Error()
			assert.Regexp(t, regexp.MustCompile(v.expected), actual)
		})
	}
}

func TestErrorLogValue(t *testing.T) {
	err := errorx.NewErrorForTest("error!", io.EOF, "at xxx", slog.LevelInfo, 500)
	expected := map[string]slog.Value{
		"msg":   slog.StringValue("error!"),
		"cause": slog.StringValue("EOF"),
		"at":    slog.StringValue("at xxx"),
	}
	for _, actualAttr := range err.LogValue().Group() {
		expectedValue, ok := expected[actualAttr.Key]
		assert.True(t, ok)
		assert.True(t, actualAttr.Value.Equal(expectedValue))
	}
}

func TestErrorHTTPStatus(t *testing.T) {
	tests := map[string]struct {
		inStatus []errorx.OptionFunc
		expected int
	}{
		"default 500": {expected: 500},
		"with status": {inStatus: []errorx.OptionFunc{errorx.WithStatus(400)}, expected: 400},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			actual := errorx.NewError("", v.inStatus...).HTTPStatus()
			assert.Equal(t, v.expected, actual)
		})
	}
}

func TestNRAttributer(t *testing.T) {
	type input struct {
		msg   string
		cause error
	}
	type want struct {
		msg   string
		cause error
	}
	tests := map[string]struct {
		in   input
		want want
	}{
		"cause is nil": {
			in: input{
				msg: "error",
			},
			want: want{
				msg: "error",
			},
		},
		"cause is not nil": {
			in: input{
				msg:   "error",
				cause: io.EOF,
			},
			want: want{
				msg:   "error",
				cause: io.EOF,
			},
		},
	}
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			err := errorx.NewError(v.in.msg, errorx.WithCause(v.in.cause))
			attr := err.NRAttribute()

			got, ok := attr["msg"]
			assert.True(t, ok)
			assert.Equal(t, v.want.msg, got)

			if v.in.cause != nil {
				got, ok = attr["cause"]
				assert.True(t, ok)
				assert.Equal(t, v.want.cause.Error(), got)
			}

			got, ok = attr["at"]
			assert.True(t, ok)
			assert.NotEmpty(t, got)
		})
	}
}
