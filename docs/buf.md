# Buf

This document provides an overview of [Buf](https://buf.build/docs/introduction) and how its capabilities can be used with Bazel.

"Buf builds tooling to make schema-driven, Protobuf-based API development reliable and user-friendly for service
producers and consumers."

## Use cases

The Buf features we're particularly interested in are:
* Accessing `.proto` files from the Buf Schema Registry, specifically:
  * [protovalidate](https://github.com/bufbuild/protovalidate)) for schema validation
  * [wellknowntypes](https://buf.build/protocolbuffers/wellknowntypes) to potentially replace those provided by [toolchains_protoc](./rules-proto.md)
* Detecting breaking changes to `.proto` files

## Known issues

Buf doesn't seem to support bzlmod effectively. Repos in the Buf Schema Registry that provide multiple packages do not
work well with Gazelle ([Issue #76](https://github.com/bufbuild/rules_buf/issues/76)).

Targets defined in `rules_buf` also does not support Protocol Buffer toolchains in Bazel ([Issue #74](https://github.com/bufbuild/rules_buf/issues/74)).

## What rules_buf provides

`rules_buf` includes a [Gazelle](./gazelle.md) extension for generating `buf_dependencies`, `buf_breaking_test`, and
`buf_lint_test` rules out of `buf.yaml` configuration files.

## buf_dependencies

[buf_dependencies](https://buf.build/docs/build-systems/bazel#buf-dependencies) is a repository rule that downloads modules from the Buf Schema Registry and generates `BUILD.bazel`
files using Gazelle.

In practice, this rule allows us to download and use `.proto` messages like those defined in [googleapis](https://buf.build/googleapis/googleapis) without
having to load the dependency in a more convoluted way.

## buf_lint_test

[buf_lint_test](https://buf.build/docs/build-systems/bazel#buf-lint-test) is a test rule that lints one or more
`proto_library` targets.

## References and links

* [rules_buf GitHub repo](https://github.com/bufbuild/rules_buf?tab=readme-ov-file)
* [Buf and Bazel compatibility overview](https://buf.build/docs/build-systems/bazel)
* [rules_buf with bzlmod example](https://github.com/bufbuild/rules_buf/blob/d18f1270a62b726e2df14c7f40186e1fe6abf602/examples/bzlmod/BUILD.bazel)
