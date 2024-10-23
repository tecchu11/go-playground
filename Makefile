.PHONY: setup, gen

setup:
	@asdf plugin add nodejs
	@asdf plugin add sqlc https://github.com/tecchu11/asdf-sqlc.git
	@asdf plugin add golangci-lint
	@asdf install
	@npm i -g @redocly/cli@latest
	@go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

gen:
	@go generate --tags oapi ./...
	@go generate --tags tools ./...
