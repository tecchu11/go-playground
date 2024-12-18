# syntax=docker/dockerfile:1

FROM golang:1.23.4 AS build
WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    GOPRIVATE=github.com/tecchu11 go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build -o /bin/migration --ldflags="-s -w" -trimpath ./cmd/migration

# latest-arm64
FROM gcr.io/distroless/static-debian12@sha256:50f827f875a7a4fc95ebbfcb309f20268065152926ff24672ec0eec70c162f21 AS final

COPY --from=build --chown=nonroot:nonroot /bin/migration /bin/

USER nonroot

ENTRYPOINT [ "/bin/migration", "up" ]
