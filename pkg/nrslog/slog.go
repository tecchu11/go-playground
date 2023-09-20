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

// OptionFunc is optional func for nrHandler.
type OptionFunc func(*Config)

// WithHandler configure specific slog.Handler to nrHandler.
func WithHandler(handler slog.Handler) OptionFunc {
	return func(config *Config) {
		config.handler = handler
	}
}

// EnableNRLogForward determines whether logs are sent via newrelic go agent.
func EnableNRLogForward() OptionFunc {
	return func(config *Config) {
		config.enableNRLogForward = true
	}
}

// nrHandler is a Handler that writes a record with newrelic metadata via the parent handler.
type nrHandler struct {
	parent             slog.Handler
	app                *newrelic.Application
	enableNRLogForward bool
	appName            string
}

// New initialize slog.Logger consisted by nrHandler with given newrelic.Application and options.
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
//		// Init NRLogger and then set Logger consited by NRHandler globally.
//		slog.SetDefault(nrslog.New(nrApp))
//	}
func New(app *newrelic.Application, opts ...OptionFunc) *slog.Logger {
	config := &Config{
		handler:            slog.Default().Handler(),
		enableNRLogForward: false,
	}
	for _, optFunc := range opts {
		optFunc(config)
	}
	nrConf, ok := app.Config()
	if !ok {
		handler := &nrHandler{
			parent:             config.handler,
			app:                app,
			enableNRLogForward: config.enableNRLogForward,
		}
		return slog.New(handler)
	}
	handler := &nrHandler{
		parent:             config.handler,
		app:                app,
		enableNRLogForward: config.enableNRLogForward,
		appName:            nrConf.AppName,
	}
	return slog.New(handler)

}

// Handle writes logs via the parent Handler with newerlic metadata from newrelic.Transaction or newrelic.Application.
// if enableNRLogForward is false, no logs are sent via newrelic go agent.
func (handler *nrHandler) Handle(ctx context.Context, record slog.Record) error {
	txn := newrelic.FromContext(ctx)
	if !handler.enableNRLogForward {
		if txn == nil {
			record.AddAttrs(slog.String(logcontext.KeyEntityName, handler.appName))
			return handler.parent.Handle(ctx, record)
		}
		record.AddAttrs(nrAttrsFromTrasnsaction(txn)...)
		return handler.parent.Handle(ctx, record)
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
	if handler.app != nil {
		handler.app.RecordLog(data)
		return nil
	}
	return handler.parent.Handle(ctx, record)
}

// Enabled reports whether the handler handles records at the given level.
// The handler ignores records whose level is lower.
func (h *nrHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.parent.Enabled(ctx, level)
}

// WithAttrs returns a new NRHandler whose attributes consists of h's attributes followed by attrs.
func (h *nrHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &nrHandler{h.parent.WithAttrs(attrs), h.app, h.enableNRLogForward, h.appName}
}

func (h *nrHandler) WithGroup(name string) slog.Handler {
	return &nrHandler{h.parent.WithGroup(name), h.app, h.enableNRLogForward, h.appName}
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
