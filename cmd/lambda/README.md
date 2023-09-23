# Lambda go

Containerized go app run on aws lambda.

## How to run locally

Specify build target `local`.

```bash
docker build -t go-playground-lambda:test . -f Dockerfile.lambda --target local

docker run -p 9000:8080 --rm go-playground-lambda:test /bin/server

curl "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{"payload":"hello world!"}'
```

## How to build for production

Specify build target `real`.

```bash
docker build -t go-playground-lambda:1.0.0 . -f Dockerfile.lambda --target real
```

