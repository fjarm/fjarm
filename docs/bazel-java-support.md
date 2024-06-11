# Java support in Bazel

This document provides an overview of how Java is configured as a basis of this Kotlin project built with Bazel.

## Java language versions

There are two relevant versions of Java that are set with configuration flags:

* The version of the source files in the repository (Controlled by `--java_language_version`)
* The version of the Java runtime that is used to execute the code and to test it (Controlled by `--java_runtime_version`)

There is a second pair of JDK and JVM used to build and execute tools, which are used in the build process, but are not
in the build results. That pair is controlled by:

* `--tool_java_language_version`
* `--tool_java_runtime_version`

## Best practices

Prefer Maven's standard directory layout (sources under `src/main/java`, tests under `src/test/java`).

Follow these guidelines when creating your `BUILD.bazel` files:

* Use one `BUILD.bazel` file per directory containing Java sources, because this improves build performance
* Every `BUILD.bazel` file should contain one java_library rule that looks like this:

```build
java_library(
    name = "directory-name",
    srcs = glob(["*.java"]),
    deps = [
        #...
    ],
)
```
* The name of the `java_library` should be the name of the directory containing the `BUILD.bazel` file
  * This makes the label of the `java_library` shorter, that is use `//package` instead of `//package:package`
* Tests should be in a matching directory under `src/test` and depend on this `java_library`

## Java toolchains

Bazel uses two types of toolchains:

* Execution toolchain is the JVM and is controlled by `--java_runtime_version`
* Compilation toolchain is the JDK and is controlled by `--java_language_version`

[This document](https://bazel.build/docs/bazel-and-java#config-java-toolchains) provides further details about configuring the JDK and JVM using local or remote sources. But for our
purposes, the default toolchains are enough.

Because the API is developed in pure Kotlin and [Android](https://developer.android.com/build/jdks#compileSdk) is the only source of restrictions for which Java version can be
used in the project, we default to Java 17.

## References and links

* [Bazel and Java overview](https://bazel.build/docs/bazel-and-java)
