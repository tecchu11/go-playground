.DEFAULT_GOAL := help

.PHONY: init

init: ## Initializes projects
	@go work init . tools
	@go install tool

.PHONY: gen,gen-oapi,gen-sqlc,gen-tbls

gen: gen-oapi gen-sqlc gen-tbls ## Generates all code.

gen-oapi: ## Generates backend code defined by open api with oapi-code-gen
	@go generate --tags oapi ./...

gen-sqlc: ## Generates backend queries with sqlc.
	@go generate --tags sqlc ./...

gen-tbls: ## Generates db schema docs with tbls.
	@go generate --tags tbls ./...

.PHONY: help
help: ## Show help.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
