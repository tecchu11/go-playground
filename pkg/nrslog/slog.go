package nrslog

import (
	"context"
	"log/slog"

	"github.com/newrelic/go-agent/v3/integrations/logcontext"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// Config contains NRHadnler behavior settings.
type Config struct {
	handler            slog.Handler
	enableNRLogForward bool
}

// NRHandlerOption is optional func for NRHandler.
type NRHandlerOption func(*Config)

// WithHandler configure specific slog.Handler to NRHandler.
func WithHandler(handler slog.Handler) NRHandlerOption {
	return func(conf *Config) {
		conf.handler = handler
	}
}

// EnableNRLogForward determines whether logs are sent via newrelic go agent.
func EnableNRLogForward() NRHandlerOption {
	return func(conf *Config) {
		conf.enableNRLogForward = true
	}
}

// NRHandler is a Handler that writes a record with newrelic metadata from the parent handler.
type NRHandler struct {
	parent             slog.Handler
	app                *newrelic.Application
	enableNRLogForward bool
	appName            string
}

// New initialize NRHandler with given newrelic.Application and options.
//
// Example:
//
//	func main() {
//		// Set JsonHanlder globally
//		slog.SetDefault(
//			slog.New(slog.NewJsonHandler(os.Stderr, &slog.HandlerOption{Level: slog.LevelInfo})),
//		)
//		// Init newrelic.Application
//		app, _ := newrelic.NewApplication(newrelic.ConfigFromEnvironment())
//		// Init NRHandler and then set Logger consited by NRHandler globally.
//		slog.SetDefault(slog.New(nrslog.New(nrApp)))
//	}
func New(app *newrelic.Application, opts ...NRHandlerOption) *NRHandler {
	conf := &Config{
		handler:            slog.Default().Handler(),
		enableNRLogForward: false,
	}
	for _, opt := range opts {
		opt(conf)
	}
	nrConf, ok := app.Config()
	if !ok {
		return &NRHandler{
			parent:             conf.handler,
			app:                app,
			enableNRLogForward: conf.enableNRLogForward,
		}
	}
	return &NRHandler{
		parent:             conf.handler,
		app:                app,
		enableNRLogForward: conf.enableNRLogForward,
		appName:            nrConf.AppName,
	}

}

// Handle writes logs via the parent Handler with newerlic metadata from newrelic.Transaction or newrelic.Application.
// if enableNRLogForward is false, no logs are sent via newrelic go agent.
func (h *NRHandler) Handle(ctx context.Context, record slog.Record) error {
	txn := newrelic.FromContext(ctx)
	if !h.enableNRLogForward {
		if txn == nil {
			record.AddAttrs(slog.String(logcontext.KeyEntityName, h.appName))
			return h.parent.Handle(ctx, record)
		}
		record.AddAttrs(nrAttrsFromTrasnsaction(txn)...)
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
	return &NRHandler{h.parent.WithAttrs(attrs), h.app, h.enableNRLogForward, h.appName}
}

func (h *NRHandler) WithGroup(name string) slog.Handler {
	return &NRHandler{h.parent.WithGroup(name), h.app, h.enableNRLogForward, h.appName}
}

func nrAttrsFromTrasnsaction(txn *newrelic.Transaction) []slog.Attr {
	return []slog.Attr{
		slog.String(logcontext.KeyTraceID, txn.GetLinkingMetadata().TraceID),
		slog.String(logcontext.KeySpanID, txn.GetLinkingMetadata().SpanID),
		slog.String(logcontext.KeyEntityName, txn.GetLinkingMetadata().EntityName),
		slog.String(logcontext.KeyEntityType, txn.GetLinkingMetadata().EntityType),
		slog.String(logcontext.KeyEntityGUID, txn.GetLinkingMetadata().EntityGUID),
		slog.String(logcontext.KeyHostname, txn.GetLinkingMetadata().Hostname),
	}
}
