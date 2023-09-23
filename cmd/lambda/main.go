package main

import (
	"context"
	"go-playground/pkg/sloglambda"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// apiGatewayProxyHandler is type func handles lambda APIGateway handler.
type apiGatewayProxyHandler func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

// slogMiddleware set default slog.Logger
func slogMiddleware(next apiGatewayProxyHandler) apiGatewayProxyHandler {
	return func(ctx context.Context, apr events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		slog.SetDefault(sloglambda.New(ctx))
		return next(ctx, apr)
	}
}

// handler handle event source APIGatewayProxyRequest.
func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	slog.InfoContext(ctx, "Request was recived")
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "\"Hello from Lambda!\"",
	}
	return response, nil
}

func main() {
	lambda.Start(slogMiddleware(handler))
}
