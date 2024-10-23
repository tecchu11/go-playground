#bin/sh

redocly bundle openapi_base.yml -o openapi.yml

oapi-codegen --config=./config.yml openapi.yml
