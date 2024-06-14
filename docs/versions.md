# Dependency versions

In addition to the automatic version tracking done by Bazel modules, this document manually tracks each of the project's
dependencies and versions in more detail.

## API

| Name | Version                                                  | Notes                            |
| ---- |----------------------------------------------------------|----------------------------------|
| rules_kotlin | [1.9.5](https://registry.bazel.build/modules/rules_kotlin) | Used to compile API source code that is written in pure Kotlin |
| rules_jvm_external | [6.1](https://registry.bazel.build/modules/rules_jvm_external) | Used to manage Maven dependencies |
| rules_java | [7.6.1](https://registry.bazel.build/modules/rules_java) | Used to configure Java toolchain |

## IDL

| Name | Version                                                         | Notes                                                          |
| ---- |-----------------------------------------------------------------|----------------------------------------------------------------|
| rules_proto | [6.0.0](https://registry.bazel.build/modules/rules_proto)       | N/A                                                            |
| toolchains_protoc | [0.3.0](https://registry.bazel.build/modules/toolchains_protoc) | Provides the `protoc` `v27.0` binary                           |
| rules_buf | [0.3.0](https://registry.bazel.build/modules/rules_buf) | Only used to pull schemas from Buf Schema Registry |
