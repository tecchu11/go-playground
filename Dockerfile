# syntax=docker/dockerfile:1

FROM golang:1.22.3 AS build
WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build -o /bin/server ./cmd/api

FROM gcr.io/distroless/static-debian12@sha256:a216254b8f42a015e380cebb6538488bed896a7072dcac951007d14f79806b84 AS final

COPY --from=build /bin/server /bin/

EXPOSE 8080

ENTRYPOINT [ "/bin/server" ]
