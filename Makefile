.PHONY: init, gen

init:
	@go work init . tools

gen:
	@go generate --tags oapi ./...
	@go generate --tags sqlc ./...
