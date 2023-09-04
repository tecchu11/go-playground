package nrslog

import (
	"context"
	"log/slog"

	"github.com/newrelic/go-agent/v3/integrations/logcontext"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// NRHandler is a Handler that writes a record with newrelic metadata from the parent handler.
type NRHandler struct {
	parent        slog.Handler
	app           *newrelic.Application
	enableConsume bool
}

// NewNRHandler creates a NRHandler.
// If you don't want to send logs to newrelic from go agent, please specify false enableConsume as false.
func NewNRHandler(parent slog.Handler, app *newrelic.Application, enableConsume bool) *NRHandler {
	return &NRHandler{
		parent:        parent,
		app:           app,
		enableConsume: enableConsume,
	}
}

// Handle writes logs by the parent Handler with newerlic metadata from newrelic.Transaction.
// If enableConsume is false, newrelic go agent is not used to log.
func (h *NRHandler) Handle(ctx context.Context, record slog.Record) error {
	txn := newrelic.FromContext(ctx)
	if !h.enableConsume {
		if txn == nil {
			return h.parent.Handle(ctx, record)
		}
		record.AddAttrs(nrAttrs(txn)...)
		return h.parent.Handle(ctx, record)
	}
	data := newrelic.LogData{
		Timestamp: record.Time.UnixMilli(),
		Severity:  record.Message,
		Message:   record.Message,
	}
	if txn != nil {
		txn.RecordLog(data)
		return nil
	}
	if h.app != nil {
		h.app.RecordLog(data)
		return nil
	}
	return h.parent.Handle(ctx, record)
}

// Enabled reports whether the handler handles records at the given level.
// The handler ignores records whose level is lower.
func (h *NRHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.parent.Enabled(ctx, level)
}

// WithAttrs returns a new NRHandler whose attributes consists of h's attributes followed by attrs.
func (h *NRHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &NRHandler{h.parent.WithAttrs(attrs), h.app, h.enableConsume}
}

func (h *NRHandler) WithGroup(name string) slog.Handler {
	return &NRHandler{h.parent.WithGroup(name), h.app, h.enableConsume}
}

func nrAttrs(txn *newrelic.Transaction) []slog.Attr {
	return []slog.Attr{
		slog.String(logcontext.KeyTraceID, txn.GetLinkingMetadata().TraceID),
		slog.String(logcontext.KeySpanID, txn.GetLinkingMetadata().SpanID),
		slog.String(logcontext.KeyEntityName, txn.GetLinkingMetadata().EntityName),
		slog.String(logcontext.KeyEntityType, txn.GetLinkingMetadata().EntityType),
		slog.String(logcontext.KeyEntityGUID, txn.GetLinkingMetadata().EntityGUID),
		slog.String(logcontext.KeyHostname, txn.GetLinkingMetadata().Hostname),
	}
}

