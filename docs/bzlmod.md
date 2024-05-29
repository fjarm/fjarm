# Bzlmod

This document provides an overview of Bazel Modules, the successor to `WORKSPACE` files that traditionally managed
external dependencies when using the Bazel build system.

"A Bazel module is a Bazel project that can have multiple versions, each of which publishes metadata about other modules
that it depends on. This is analogous to familiar concepts in other dependency management systems, such as a Maven
artifact, an npm package, a Go module, or a Cargo crate."

## Bazel lockfile

The lockfile feature in Bazel records specific versions or dependencies of libraries or packages required by a project.
It enhances build efficiency by allowing Bazel to skip the resolution process when there are no changes in project
dependencies.

The `MODULE.bazel.lock` file is created or updated during the build process, specifically after module resolution and
extension evaluation. Importantly, it only includes dependencies that are included in the current invocation of the
build.

The lockfile can be controlled by the flag `--lockfile_mode` to customize the behavior of Bazel when the project state
differs from the lockfile. The available modes are:

* `update` (Default): If the project state matches the lockfile, the resolution result is immediately returned from the 
lockfile. Otherwise, resolution is executed, and the lockfile is updated to reflect the current state.
* `error`: If the project state matches the lockfile, the resolution result is returned from the lockfile. Otherwise,
Bazel throws an error indicating the variations between the project and the lockfile. This mode is particularly useful
when you want to ensure that your project's dependencies remain unchanged, and any differences are treated as errors.
* `off`: The lockfile is not checked at all.

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
