//go:build tools

// Package tools provides all tolls for this app.
//
// tool.go exists so that all tools can be managed by Renovate.
package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
	_ "github.com/k1LoW/tbls"
)
