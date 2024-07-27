package nrslog

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/newrelic/go-agent/v3/integrations/logcontext"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// nrJSONHandler is a Handler that writes a record with newrelic metadata via the parent handler.
type nrJSONHandler struct {
	slog.Handler
	app *newrelic.Application
}

// NewJSONHandler creates nrJSONHandler.
func NewJSONHandler(app *newrelic.Application, opts *slog.HandlerOptions) (slog.Handler, error) {
	base := slog.NewJSONHandler(os.Stdout, opts)
	conf, ok := app.Config()
	if !ok {
		return nil, errors.New("missing newrelic.Application because of Application being not yet fully initialized")
	}
	decorated := base.WithAttrs([]slog.Attr{slog.String(logcontext.KeyEntityName, conf.AppName)})
	return &nrJSONHandler{Handler: decorated, app: app}, nil
}

// Handle writes logs via the parent Handler with newrelic metadata from newrelic.Transaction or newrelic.Application.
func (h *nrJSONHandler) Handle(ctx context.Context, record slog.Record) error {
	txn := newrelic.FromContext(ctx)
	if txn == nil {
		return h.Handler.Handle(ctx, record)
	}
	record.AddAttrs(nrAttrs(txn)...)
	return h.Handler.Handle(ctx, record)
}

// WithAttrs returns a new nrJSONHandler whose attributes consists of h's attributes followed by attrs.
func (h *nrJSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &nrJSONHandler{Handler: h.Handler.WithAttrs(attrs), app: h.app}
}

func (h *nrJSONHandler) WithGroup(name string) slog.Handler {
	return &nrJSONHandler{Handler: h.Handler.WithGroup(name), app: h.app}
}

func nrAttrs(txn *newrelic.Transaction) []slog.Attr {
	return []slog.Attr{
		slog.String(logcontext.KeyTraceID, txn.GetLinkingMetadata().TraceID),
		slog.String(logcontext.KeySpanID, txn.GetLinkingMetadata().SpanID),
		slog.String(logcontext.KeyEntityType, txn.GetLinkingMetadata().EntityType),
		slog.String(logcontext.KeyEntityGUID, txn.GetLinkingMetadata().EntityGUID),
		slog.String(logcontext.KeyHostname, txn.GetLinkingMetadata().Hostname),
	}
}

var _ slog.Handler = (*nrJSONHandler)(nil)
