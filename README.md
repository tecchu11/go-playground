# Go Playground

go playground for tecchu11

## How to

### Build and run image

build

```bash
GITHUB_TOKEN=$(gh auth token) docker build \
--secret id=GITHUB_TOKEN \
-t go-playground:latest . 
```

build(debug)

```bash
GITHUB_TOKEN=$(gh auth token) docker build \
--secret id=GITHUB_TOKEN \
--target debug \
-t go-playground:latest . 
```

run
```bash
docker run -p 8080:8080 \
--name go-playground \
--env-file .env go-playground:latest
```

exec(debug)

```bash
docker exec -it go-playground sh
```

### Debug renovate config

Check regex pattern with [regex101](https://regex101.com/)

```bash
TOKEN=$(gh auth token) && \
    RENOVATE_BASE_BRANCHES=xxx \
    LOG_LEVEL=debug \
    RENOVATE_CONFIG_FILE=renovate.json \
    renovate \
    --token "$TOKEN" \
    --dry-run \
    tecchu11/go-playground
```
