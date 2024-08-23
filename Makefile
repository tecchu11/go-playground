.PHONY: setup

setup:
	@asdf plugin add nodejs
	@asdf plugin add sqlc https://github.com/tecchu11/asdf-sqlc.git
	@asdf plugin add golangci-lint
	@asdf install
