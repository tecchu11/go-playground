#bin/sh

BASE_PATH=../../../../../../doc/oas

sh $BASE_PATH/bundle.sh
oapi-codegen --config=./config.yml $BASE_PATH/openapi.yml
