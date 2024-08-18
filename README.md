# Go Playground

go playground for tecchu11

## Debug renovate config

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
