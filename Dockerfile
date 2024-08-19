# syntax=docker/dockerfile:1

FROM golang:1.23.0 AS build
WORKDIR /src

ENV GOPRIVATE=github.com/tecchu11

RUN --mount=type=secret,id=GITHUB_TOKEN \
    git config --global \
        url."https://x-access-token:$(cat /run/secrets/GITHUB_TOKEN)@github.com/tecchu11/".insteadOf \
        https://github.com/tecchu11/

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build -o /bin/server ./cmd/api

FROM gcr.io/distroless/static-debian12:debug-nonroot-arm64 AS debug

COPY --from=build --chown=nonroot:nonroot /bin/server /bin/

EXPOSE 8080

ENTRYPOINT [ "/bin/server" ]

# nonroot-arm64
FROM gcr.io/distroless/static-debian12@sha256:53fe5a9786d457f03bdb9cdb29f06b0c543c54c02bb9e9bfdb62709aa7a09825 AS final

COPY --from=build --chown=nonroot:nonroot /bin/server /bin/

EXPOSE 8080

ENTRYPOINT [ "/bin/server" ]
