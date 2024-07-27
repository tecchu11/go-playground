package main

import (
	"fmt"
	"log/slog"
	"net/rpc"
	"os"

	"github.com/aws/aws-lambda-go/lambda/messages"
)

type lambdaClient struct {
	client *rpc.Client
}

func newClient() (*lambdaClient, error) {
	client, err := rpc.Dial("tcp", "localhost:9000")
	if err != nil {
		return nil, err
	}
	return &lambdaClient{client: client}, nil
}

func (lc *lambdaClient) ping() error {
	req := &messages.PingRequest{}
	var res *messages.PingResponse
	return lc.client.Call("Function.Ping", req, &res)
}

func (lc *lambdaClient) invoke(payload []byte) ([]byte, error) {
	req := &messages.InvokeRequest{Payload: payload}
	res := messages.InvokeResponse{}

	if err := lc.client.Call("Function.Invoke", req, &res); err != nil {
		return nil, err
	}
	return res.Payload, nil
}

func main() {
	lc, err := newClient()
	if err != nil {
		slog.Error("Dial error", slog.String("err", err.Error()))
		os.Exit(1)
	}
	if err := lc.ping(); err != nil {
		slog.Error("Ping error", slog.String("err", err.Error()))
		os.Exit(1)
	}
	res, err := lc.invoke([]byte("{\"foo\":100}"))
	if err != nil {
		slog.Error("Invoke error", slog.String("err", err.Error()))
		os.Exit(1)
	}
	fmt.Println(string(res))
}
