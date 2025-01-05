# Protocol buffers in Bazel

This document provides an overview of how to use Protocol Buffers with Bazel using the base layer provided by
`rules_proto`.

## What rules_proto provides

`rules_proto` provides two rules: `proto_library` and `proto_lang_toolchain`.

`proto_library` defines libraries of protocol buffers which may be used from multiple languages. A `proto_library` may
be listed in the `deps` clause of supported rules, such as `java_proto_library`.

When compiled on the command-line, a `proto_library` creates a file named `foo-descriptor-set.proto.bin`, which is the
descriptor set for the messages in the rule `srcs`. The file is a serialized `FileDescriptorSet`, which is described
[here](https://developers.google.com/protocol-buffers/docs/techniques#self-description).

It only contains information about the `.proto` files directly mentioned by a `proto_library` rule.

`proto_lang_toolchain` will be covered later.

## rules_proto best practices

Recommended code organization:

* One `proto_library` rule per `.proto` file.
* A file named `foo.proto` will be in a rule named `foo_proto`, which is located in the same package.
* A `[language]_proto_library` that wraps a `proto_library` named `foo_proto` should be called `foo_[language]_proto`,
and be located in the same package.


## What toolchains_protoc provides

`toolchains_protoc` provides a Bazel module extension that 1) allows us to specify a versioned `protoc` binary (or
"toolchain") to use when compiling `.proto` files, and 2) provides us a repository to specify well-known types like
`@com_google_protobuf//:empty_proto`.

A pre-built, versioned `protoc` binary is useful because the alternative would be to rely on `rules_proto` downloading
`protoc` each time we build `proto_library` rule. This would either require specifying a C++ toolchain to ensure that
each build is hermetic even when the codebase doesn't have any C++ source files. Or, we would have to rely on the host
machine's C++ version, which negates Bazel's reproducible build benefit.

Also, downloading and compiling `protoc` through `rules_proto` leads to increased build times. We can verify this by
removing the `protoc` extension from `MODULE.bazel` and building `//proto/dummy/v1:dummy_proto`. Notice the cluttered
logs like:

```
$ external/protobuf~/src/google/protobuf/compiler/cpp/helpers.cc:197:25: warning: unused function 'VerifyInt32TypeToVerifyCustom' [-Wunused-function]
```

## References and links

* [rules_proto GitHub repo](https://github.com/bazelbuild/rules_proto)
* [Bazel encyclopedia proto_library documentation](https://bazel.build/reference/be/protocol-buffer#proto_library)
