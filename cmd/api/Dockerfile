# syntax=docker/dockerfile:1

FROM golang:1.25.6 AS build
WORKDIR /src

RUN --mount=type=secret,id=GITHUB_TOKEN \
    git config --global \
        url."https://x-access-token:$(cat /run/secrets/GITHUB_TOKEN)@github.com/tecchu11/".insteadOf \
        https://github.com/tecchu11/

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    GOPRIVATE=github.com/tecchu11 go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build -o /bin/server --ldflags="-s -w" -trimpath ./cmd/api

FROM gcr.io/distroless/static-debian12:debug AS debug

COPY --from=build --chown=nonroot:nonroot /bin/server /bin/

USER nonroot

EXPOSE 8080

ENTRYPOINT [ "/bin/server" ]

# latest-arm64
FROM gcr.io/distroless/static-debian12@sha256:cd64bec9cec257044ce3a8dd3620cf83b387920100332f2b041f19c4d2febf93 AS final

COPY --from=build --chown=nonroot:nonroot /bin/server /bin/

USER nonroot

EXPOSE 8080

ENTRYPOINT [ "/bin/server" ]
