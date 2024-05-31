# Dependency versions

In addition to the automatic version tracking done by Bazel modules, this document manually tracks each of the project's
dependencies and versions in more detail.

## Protocol Buffer

| Name | Version                                                               | Notes                                                          |
| ---- |-----------------------------------------------------------------------|----------------------------------------------------------------|
| rules_proto | [6.0.0](https://registry.bazel.build/modules/rules_proto/6.0.0)       | N/A                                                            |
| toolchains_protoc | [0.3.0](https://registry.bazel.build/modules/toolchains_protoc/0.3.0) | Provides the `protoc` `v27.0` binary                           |
| rules_buf | [0.3.0](https://registry.bazel.build/modules/rules_buf/0.3.0) | Only used to pull schemas from Buf Schema Registry |
