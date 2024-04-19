package slogx

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

var isColdStart = true

// NewLogger initializes slog.Logger which configured lambda context.
func NewLogger(ctx context.Context, opt *slog.HandlerOptions) *slog.Logger {
	if opt == nil {
		opt = &slog.HandlerOptions{}
	}
	handler := slog.NewJSONHandler(os.Stdout, opt)
	var coldStart bool
	if isColdStart {
		coldStart = true
		isColdStart = false
	}
	lambdaCtx, ok := lambdacontext.FromContext(ctx)
	if !ok {
		decorated := handler.WithAttrs([]slog.Attr{
			slog.Bool("coldStart", coldStart),
			slog.Group("function",
				slog.String("name", lambdacontext.FunctionName),
				slog.String("version", lambdacontext.FunctionVersion)),
		})
		return slog.New(decorated)
	}
	decorated := handler.WithAttrs([]slog.Attr{
		slog.Bool("coldStart", coldStart),
		slog.Group("function",
			slog.String("name", lambdacontext.FunctionName),
			slog.String("version", lambdacontext.FunctionVersion)),
		slog.String("requestID", lambdaCtx.AwsRequestID),
	})
	return slog.New(decorated)
}
