.PHONY: fmt, test, gen, migration-up, migration-down

fmt:
	@go fmt ./...

test:
	@go test ./... --race

gen:
	@go generate -tags tools  ./...    

migration-up:
	@go run ./cmd/migration up .env

migration-down:
	@go run ./cmd/migration down .env
