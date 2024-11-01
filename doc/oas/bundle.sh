#bin/sh

REDOCLY_TAG=1.25.10
BASE="$(cd "$(dirname "$0")" && pwd)"

docker run --rm -v $BASE:/spec --entrypoint sh redocly/cli:$REDOCLY_TAG -c "
    redocly bundle /spec/openapi_base.yml -o /spec/openapi.yml &&
    redocly lint /spec/openapi.yml \
        --skip-rule operation-4xx-response \
        --skip-rule no-server-example.com \
        --skip-rule info-license \
        --skip-rule security-defined
"
