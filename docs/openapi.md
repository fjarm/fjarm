# OpenAPI

## Resources and links

* Generate OpenAPI Kotlin models and a client with
```bash
docker run --rm -v "${PWD}:/local" openapitools/openapi-generator-cli generate -i /local/openapi/dummy/v1/dummy.yaml -i /local/openapi/dummy/v1/dummy_service.yaml -g kotlin -o /local/out/ --model-package xyz.fjarm.dummy.v1 --ignore-file-override /local/openapi/dummy/v1/.openapi-generator-ignore -c /local/openapi/dummy/v1/config.yaml
```