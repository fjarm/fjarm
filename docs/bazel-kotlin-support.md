# Kotlin support in Bazel

This is an overview of how to use Kotlin with Bazel as a build system. For more details, refer to [`rules_kotlin`](https://github.com/bazelbuild/rules_kotlin).

## rules_jvm_external dependency

[`rules_jvm_external`](https://registry.bazel.build/modules/rules_jvm_external) is depended on to make managing Maven dependencies like the
Protobuf Java runtime.

## rules_java dependency

[`rules_java`](https://registry.bazel.build/modules/rules_java) can be added as a direct dependency in the project in order to configure Java toolchains.

This is needed because we need to compile `proto_library` targets into `java_proto_library` targets that can be used in
Kotlin source code.
