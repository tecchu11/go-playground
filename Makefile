.PHONY: fmt, test, gen

fmt:
	@go fmt ./...

test:
	@go test ./... --race

gen:
	@go generate -tags tools  ./...    
