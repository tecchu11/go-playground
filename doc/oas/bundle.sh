#bin/sh

BASE="$(cd "$(dirname "$0")" && pwd)"

redocly bundle $BASE/openapi_base.yml -o $BASE/openapi.yml
redocly lint $BASE/openapi.yml \
    --skip-rule operation-4xx-response \
    --skip-rule no-server-example.com \
    --skip-rule info-license \
    --skip-rule security-defined
