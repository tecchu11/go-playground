# syntax=docker/dockerfile:1

FROM ghcr.io/k1low/tbls:v1.79.3 AS build
WORKDIR /src

RUN --mount=type=bind,source=.tbls.yml,target=.tbls.yml \
    tbls doc

FROM squidfunk/mkdocs-material:9.5  AS final

COPY --from=build /src/dest/ /docs
COPY mkdocs.yml .

RUN pip install mkdocs-glightbox

EXPOSE 8000

ENTRYPOINT ["/sbin/tini", "--", "mkdocs"]
CMD ["serve", "--dev-addr=0.0.0.0:8000"]
