.PHONY: fmt, test

fmt:
	@go fmt ./...

test:
	@go test ./... --race

update-all:
	@go get -u ./...
	@go mod tidy
