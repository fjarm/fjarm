# Bzlmod

This document provides an overview of Bazel Modules, the successor to `WORKSPACE` files that traditionally managed
external dependencies when using the Bazel build system.

"A Bazel module is a Bazel project that can have multiple versions, each of which publishes metadata about other modules
that it depends on. This is analogous to familiar concepts in other dependency management systems, such as a Maven
artifact, an npm package, a Go module, or a Cargo crate."

## Module tags and extensions

Modules can also specify customized pieces of data called tags, which are consumed by module extensions after module
resolution to define additional repos.

These extensions have capabilities similar to repo rules, enabling them to perform actions like file I/O and sending
network requests. Among other things, they allow Bazel to interact with other package management systems while also
respecting the dependency graph built out of Bazel modules.

## Resources and links

* [Bazel external dependencies overview](https://bazel.build/external/overview)
* [Bazel modules user guide](https://bazel.build/external/module)
* [Full list of directives available in MODULE.bazel](https://bazel.build/rules/lib/globals/module)
* [Bazel central registry](https://registry.bazel.build/)
* [Bazel lockfile overview](https://bazel.build/external/lockfile)
