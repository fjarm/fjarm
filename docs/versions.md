# Dependency versions

In addition to the automatic version tracking done by Bazel modules, this document manually tracks each of the project's
dependencies and versions in more detail.

## API

| Name | Version                                                  | Notes                            |
| ---- |----------------------------------------------------------|----------------------------------|
| | |

## IDL

| Name | Version                                                         | Notes                                                              |
| ---- |-----------------------------------------------------------------|--------------------------------------------------------------------|
| rules_proto | [6.0.2](https://registry.bazel.build/modules/rules_proto)       | N/A                                                                |
| toolchains_protoc | [0.3.1](https://registry.bazel.build/modules/toolchains_protoc) | Provides the `protoc` `v27.0` binary                               |
| rules_buf | [0.3.0](https://registry.bazel.build/modules/rules_buf)         | Used to pull schemas from Buf Schema Registry using `buf` `v1.34.0` |
