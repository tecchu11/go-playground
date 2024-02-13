# syntax=docker/dockerfile:1

FROM golang:1.22.0 AS build
WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build -o /bin/server ./cmd/api

FROM gcr.io/distroless/static-debian12@sha256:711bbee1fca4bac059794b8887fb6da6a98941d4018c94c5a3c527e85d77a811 AS final

COPY --from=build /bin/server /bin/

EXPOSE 8080

ENTRYPOINT [ "/bin/server" ]
