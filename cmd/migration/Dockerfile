# syntax=docker/dockerfile:1

FROM golang:1.24.1 AS build
WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    GOPRIVATE=github.com/tecchu11 go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build -o /bin/migration --ldflags="-s -w" -trimpath ./cmd/migration

# latest-arm64
FROM gcr.io/distroless/static-debian12@sha256:49f2d1ff9571823cc6470dd4f068b3b8c3b6552a06db5204e93be0ef5c6f3bdd AS final

COPY --from=build --chown=nonroot:nonroot /bin/migration /bin/

USER nonroot

ENTRYPOINT [ "/bin/migration", "up" ]
