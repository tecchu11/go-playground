.PHONY: setup, fmt, test, gen, migration-up, migration-down

setup:
	@asdf plugin add nodejs
	@asdf plugin add sqlc https://github.com/tecchu11/asdf-sqlc.git
	@asdf plugin add golangci-lint
	@asdf install

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
