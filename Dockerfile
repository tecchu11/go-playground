# syntax=docker/dockerfile:1

FROM golang:1.23.1 AS build
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
    CGO_ENABLED=0 go build -o /bin/server ./cmd/api

FROM gcr.io/distroless/static-debian12:debug AS debug

COPY --from=build --chown=nonroot:nonroot /bin/server /bin/

USER nonroot

EXPOSE 8080

ENTRYPOINT [ "/bin/server" ]

# latest-arm64
FROM gcr.io/distroless/static-debian12@sha256:31034c10eadd77e11256b49b8fcb9d6009daa01efbd76f078bd160249fb823d6 AS final

COPY --from=build --chown=nonroot:nonroot /bin/server /bin/

USER nonroot

EXPOSE 8080

ENTRYPOINT [ "/bin/server" ]
