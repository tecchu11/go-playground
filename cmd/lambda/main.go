package main

import (
	"context"
	"go-playground/pkg/lambda/slogx"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// handler handle event source APIGatewayProxyRequest.
func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	slog.InfoContext(ctx, "Request was received")
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "\"Hello from Lambda!\"",
	}
	return response, nil
}

type hn func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

var middleware = func(next hn) hn {
	return func(ctx context.Context, apr events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		log := slogx.NewLogger(ctx, nil)
		slog.SetDefault(log)
		return next(ctx, apr)
	}
}

func main() {
	lambda.Start(middleware(handler))
}
