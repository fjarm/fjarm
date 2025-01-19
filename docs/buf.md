# Buf

This document provides an overview of [Buf](https://buf.build/docs/introduction) and how its capabilities can be used to
manage Protocol Buffer definitions defined here and upstream Protocol Buffer dependencies.

"Buf builds tooling to make schema-driven, Protobuf-based API development reliable and user-friendly for service
producers and consumers."

## What Buf provides

## buf lint

`buf lint` can be used from the project's root directory to run lint tests on schema and RPC definitions.

## [buf breaking](https://buf.build/docs/breaking/tutorial)

`buf breaking --against 'https://github.com/fjarm/fjarm.git#branch=main,subdir=proto'` is used to check for breaking
changes in `.proto` modifications.

## buf dep update

`buf dep update` is used to update the `deps` section of [buf.yaml](../buf.yaml), which specifies the external,
imported `.proto` files.

This includes `buf/validate/validate.proto`.

To introduce new dependencies:
1. Add the dependency to `buf.yaml` under the `deps` section:
```yaml
deps:
  - buf.build/googleapis/googleapis:e93e34f48be043dab55be31b4b47f458
```
2. From the project root, run `buf dep update`
3. This should result in changes to `buf.lock`

## buf generate

`buf generate` generates the language Protobuf messages and RPCs specified in a `buf.gen.yaml` file.

Since migrating to Buf Schema Registry's generated SDKs, this command is no longer used in any standard workflow.

The single exception is to validate that the upstream SDK generation works correctly.

## buf build

`buf build [DIRECTORY|MODULE] --output ./[DIRECTORY|MODULE]/image.binpb` builds the `.proto` files found in the
specified module into a Buf image. Similar to `buf generate`, this command is no longer used except to locally validate
that the upstream SDK generation will work correctly.

## buf push

`buf push [DIRECTORY|MODULE]` pushes the image associated with the module up to the Buf schema registry. This should
only be executed from the GitHub Actions UI. 

## References and links

* [Buf style guide](https://buf.build/docs/best-practices/style-guide)
* [Buf files and packages](https://buf.build/docs/reference/protobuf-files-and-packages)
* [Buf CLI docs](https://buf.build/docs/reference/cli/buf/)
* [Buf Schema Registry docs](https://buf.build/docs/bsr/introduction)
* [Buf GitHub Actions](https://buf.build/docs/ci-cd/github-actions)
