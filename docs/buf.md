# Buf

This document provides an overview of [Buf](https://buf.build/docs/introduction) and how its capabilities can be used with Bazel.

"Buf builds tooling to make schema-driven, Protobuf-based API development reliable and user-friendly for service
producers and consumers."

## What rules_buf provides

`rules_buf` includes a [Gazelle](./gazelle.md) extension for generating `buf_dependencies`, `buf_breaking_test`, and
`buf_lint_test` rules out of `buf.yaml` configuration files.

## buf_dependencies

[buf_dependencies](https://buf.build/docs/build-systems/bazel#buf-dependencies) is a repository rule that downloads modules from the Buf Schema Registry and generates `BUILD.bazel`
files using Gazelle.

In practice, this rule allows us to download and use `.proto` messages like those defined in [googleapis](https://buf.build/googleapis/googleapis) without
having to load the dependency in a more convoluted way.

## What buf provides

## buf lint

`buf lint` can be used from the [proto](../proto) folder to run lint tests on schema and RPC definitions.

## buf dep update

`buf dep update` is used to update the `deps` section of [buf.yaml](../proto/buf.yaml), which specifies the external,
imported `.proto` files.

This includes `buf/validate/validate.proto`.

## buf generate

`buf generate` generates the language Protobuf messages and RPCs specified in [buf.gen.yaml](../proto/buf.gen.yaml).

## buf build

`buf build [DIRECTORY|MODULE]` builds the `.proto` files found in the specified module into a Buf image.

## buf push

`buf push [DIRECTORY|MODULE]` pushes the image associated with the module up to the Buf schema registry. 

## References and links

* [rules_buf GitHub repo](https://github.com/bufbuild/rules_buf?tab=readme-ov-file)
* [Buf style guide](https://buf.build/docs/best-practices/style-guide)
* [Buf files and packages](https://buf.build/docs/reference/protobuf-files-and-packages)
* [Buf CLI docs](https://buf.build/docs/reference/cli/buf/)
* [Buf Schema Registry docs](https://buf.build/docs/bsr/introduction)

---

# Graveyard

## Use cases

The Buf features we're particularly interested in are:
* Accessing `.proto` files from the Buf Schema Registry, specifically:
  * [protovalidate](https://github.com/bufbuild/protovalidate)) for schema validation
  * [wellknowntypes](https://buf.build/protocolbuffers/wellknowntypes) to potentially replace those provided by [toolchains_protoc](./rules-proto.md)
* Linting `.proto` files
* Detecting breaking changes to `.proto` files

## Known issues

Buf doesn't seem to support bzlmod effectively. Repos in the Buf Schema Registry that provide multiple packages do not
work well with Gazelle ([Issue #76](https://github.com/bufbuild/rules_buf/issues/76)).

Targets defined in `rules_buf` also does not support Protocol Buffer toolchains in Bazel ([Issue #74](https://github.com/bufbuild/rules_buf/issues/74)).
