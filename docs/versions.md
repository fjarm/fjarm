# Dependency versions

In addition to the automatic version tracking done by Bazel modules, this document manually tracks each of the project's
dependencies and versions in more detail.

## Containers

| Name      | Version                                                  | Notes                                                       |
|-----------|----------------------------------------------------------|-------------------------------------------------------------|
| rules_oci | [2.0.0](https://registry.bazel.build/modules/rules_oci)  | N/A                                                         |
| rules_pkg | [0.10.0](https://registry.bazel.build/modules/rules_pkg) | Tars `go_binary` rules for eventual layering into container |

## API

| Name      | Version                                                  | Notes                                              |
|-----------|----------------------------------------------------------|----------------------------------------------------|
| rules_go  | [0.48.1](https://registry.bazel.build/modules/rules_go)  | N/A                                                |
| gazelle   | [0.37.0](https://registry.bazel.build/modules/gazelle)   | Used to automatically generate `BUILD.bazel` files |

## IDL

| Name              | Version                                                         | Notes                                                               |
|-------------------|-----------------------------------------------------------------|---------------------------------------------------------------------|
| rules_proto       | [6.0.2](https://registry.bazel.build/modules/rules_proto)       | N/A                                                                 |
| rules_buf         | [0.3.0](https://registry.bazel.build/modules/rules_buf)         | Used to pull schemas from Buf Schema Registry using `buf` `v1.34.0` |
| toolchains_protoc | [0.3.1](https://registry.bazel.build/modules/toolchains_protoc) | Provides the `protoc` `v27.0` binary                                |
