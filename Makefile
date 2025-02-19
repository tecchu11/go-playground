##@ All shortcut collections(generates code, install tools, etc...) for app script. 

##@
##@ [Set up commands]
##@

setup: ##@  Set up develop. Install all tools
	-@go work init . tools
	@make --no-print-directory install-tools

##@
##@ [Generates files commands]
##@

gen: gen-oapi gen-sqlc gen-tbls ##@ Generates all files(oapi, sqlc and tbls).

gen-oapi: ##@ Generates backend code defined by open api with oapi-code-gen
	@go generate --tags oapi ./...

gen-sqlc: ##@ Generates backend queries with sqlc.
	@go generate --tags sqlc ./...

gen-tbls: ##@ Generates db schema docs with tbls.
	@go generate --tags tbls ./...

##@
##@ [Install tools commands]
##@

install-tools: install-tools-go install-tools-tbls install-tools-redoc install-tools-renovate ##@ Install all tools for this app.

install-tools-go: ##@ Install tools via go install tool.
	@go install tool

TBLS_VERSION=1.81.0
install-tools-tbls: ##@ Install tbls.
	@curl -o tbls.deb -L https://github.com/k1LoW/tbls/releases/download/v$(TBLS_VERSION)/tbls_$(TBLS_VERSION)-1_arm64.deb
	@sudo dpkg -i tbls.deb
	@rm tbls.deb

install-tools-redoc: ##@ Install redoc.
	@npm i -g @redocly/cli@1.29.0

install-tools-renovate: ##@ Install renovate cli.
	@npm i -g renovate

##@
##@ [Misc commands]
##@

help: ##@ (Default) Show helps.
	@printf "\nUsage: make <command>\n"
	@grep -F -h "##@" $(MAKEFILE_LIST) | grep -F -v grep -F | sed -e 's/\\$$//' | awk 'BEGIN {FS = ":*[[:space:]]*##@[[:space:]]*"}; \
	{ \
		if ($$2 == "") \
			next; \
		else if ($$0 ~ /^#/) \
			printf "\n%s\n", $$2; \
		else if ($$1 ~ / /) { \
			split($$1, targets, " "); \
			if (!seen[targets[1]]++) \
				printf "\n    \033[34m%-20s\033[0m %s\n", targets[1], $$2; \
		} else \
			printf "\n    \033[34m%-20s\033[0m %s\n", $$1, $$2; \
	}'
.DEFAULT_GOAL := help
MAKEFLAGS += --no-print-directory
