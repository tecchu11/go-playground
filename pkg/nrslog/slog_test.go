package nrslog_test

import (
	"bytes"
	"context"
	"encoding/json"
	"go-playground/pkg/nrslog"
	"log/slog"
	"testing"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/require"
)

func TestNewHandler_Success(t *testing.T) {
	_, _, h := setup(t)

	require.NotNil(t, h)
}

func TestNewHandler_AppIsNotValid(t *testing.T) {
	h, err := nrslog.NewHandler(&newrelic.Application{}, slog.Default().Handler())

	require.Error(t, err)
	require.Zero(t, h)
}

func TestHandler_WithoutContext(t *testing.T) {
	_, buf, h := setup(t)
	logger := slog.New(h)
	logger.Info("test")

	type Record struct {
		Name string `json:"entity.name"`
	}
	var reocrd Record
	err := json.Unmarshal(buf.Bytes(), &reocrd)
	require.NoError(t, err)
	require.Equal(t, "test-app", reocrd.Name)
}

func TestHandler_WithContext(t *testing.T) {
	app, buf, h := setup(t)
	logger := slog.New(h)

	txn := app.StartTransaction("test")
	defer txn.End()
	ctx := newrelic.NewContext(context.Background(), txn)
	logger.InfoContext(ctx, "test")

	type Record struct {
		Name     string  `json:"entity.name"`
		TraceID  *string `json:"trace.id"`
		SpanID   *string `json:"span.id"`
		Type     *string `json:"entity.type"`
		GUID     *string `json:"entity.guid"`
		Hostname *string `json:"hostname"`
	}
	var record Record
	err := json.Unmarshal(buf.Bytes(), &record)

	require.NoError(t, err)
	require.Equal(t, "test-app", record.Name)
	require.NotNil(t, record.TraceID)
	require.NotNil(t, record.SpanID)
	require.NotNil(t, record.Type)
	require.NotNil(t, record.GUID)
	require.NotNil(t, record.Hostname)
}

func TestHandler_WithAttr(t *testing.T) {
	_, buf, h := setup(t)
	h = h.WithAttrs([]slog.Attr{slog.String("foo", "bar")})
	logger := slog.New(h)

	logger.Info("test")

	type Record struct {
		Foo string `json:"foo"`
	}
	var record Record
	err := json.Unmarshal(buf.Bytes(), &record)

	require.NoError(t, err)
	require.Equal(t, "bar", record.Foo)
}

func TestHandler_WithGroup(t *testing.T) {
	_, buf, h := setup(t)
	logger := slog.New(h.WithGroup("group"))

	logger.Info("test", slog.String("foo", "bar"))

	type Record struct {
		Group struct {
			Foo string `json:"foo"`
		} `json:"group"`
	}
	var record Record
	err := json.Unmarshal(buf.Bytes(), &record)

	require.NoError(t, err)
	require.Equal(t, "bar", record.Group.Foo)
}

func setup(t *testing.T) (*newrelic.Application, *bytes.Buffer, slog.Handler) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigLicense("0000000000000000000000000000000000000000"),
		newrelic.ConfigAppName("test-app"),
	)
	require.NoError(t, err)
	buf := bytes.NewBuffer(nil)
	handler, err := nrslog.NewHandler(app, slog.NewJSONHandler(buf, nil))
	require.NoError(t, err)
	return app, buf, handler
}
