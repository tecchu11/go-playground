# syntax=docker/dockerfile:1

FROM golang:1.25.1 AS build
WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    GOPRIVATE=github.com/tecchu11 go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build -o /bin/migration --ldflags="-s -w" -trimpath ./cmd/migration

# latest-arm64
FROM gcr.io/distroless/static-debian12@sha256:87bce11be0af225e4ca761c40babb06d6d559f5767fbf7dc3c47f0f1a466b92c AS final

COPY --from=build --chown=nonroot:nonroot /bin/migration /bin/

USER nonroot

ENTRYPOINT [ "/bin/migration", "up" ]
