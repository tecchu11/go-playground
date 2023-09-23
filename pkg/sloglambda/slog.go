package sloglambda

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

func New(ctx context.Context) *slog.Logger {
	lc, _ := lambdacontext.FromContext(ctx)
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}).WithAttrs(
		[]slog.Attr{
			slog.String("name", lambdacontext.FunctionName),
			slog.String("version", lambdacontext.FunctionVersion),
			slog.String("request_id", lc.AwsRequestID),
		},
	)
	return slog.New(handler)
}
