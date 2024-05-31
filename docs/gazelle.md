# Gazelle

This document contains an overview of the [Gazelle](https://github.com/bazelbuild/bazel-gazelle?tab=readme-ov-file) `BUILD.bazel` file generator tool.
Gazelle is particularly useful for Bazel projects that natively support Protocol Buffers.

"Gazelle is a build file generator for Bazel projects. It can create new BUILD.bazel files for a project that follows
language conventions, and it can update existing build files to include new sources, dependencies, and options."

Notably, Gazelle is [recommended in order to use `rules_buf`](https://buf.build/docs/build-systems/bazel#gazelle).

## Known issues

Gazelle helps us interact with [Buf](./buf.md) by generating `buf_dependencies` rules and updating the `deps` in
`proto_library` targets. But, as of 2024-05-30, this seems to rely on `gazelle-update-repos` being able to write a macro
[into WORKSPACE](https://buf.build/docs/build-systems/bazel#gazelle-dependencies) which isn't supported by bzlmod
([Issue #52](https://github.com/bufbuild/rules_buf/issues/52)).

## References and links

* [Gazelle GitHub repo](https://github.com/bazelbuild/bazel-gazelle?tab=readme-ov-file)
