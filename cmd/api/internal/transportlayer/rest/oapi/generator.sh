#bin/sh

BASE_PATH=../../../../../../doc/oas

redocly bundle $BASE_PATH/openapi_base.yml -o $BASE_PATH/openapi.yml
redocly lint $BASE_PATH/openapi.yml --skip-rule operation-4xx-response --skip-rule no-server-example.com --skip-rule info-license --skip-rule security-defined

oapi-codegen --config=./config.yml $BASE_PATH/openapi.yml
