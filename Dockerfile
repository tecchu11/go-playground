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

FROM gcr.io/distroless/static-debian12@sha256:49d4bd9394ec63ecadaca2afb5520c9bbb6f6b3a144418a883a9728be6edf7a1 AS final

COPY --from=build /bin/server /bin/

EXPOSE 8080

ENTRYPOINT [ "/bin/server" ]
