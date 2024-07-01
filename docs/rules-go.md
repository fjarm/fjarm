# Go with Bazel

This doc provides an overview of how to work with Go in Bazel using [rules_go](https://github.com/bazelbuild/rules_go/).

## Using the `go` SDK with Bazel

Instead of using the `go` SDK directly from the command line then setting Bazel up to work with the results of those
`go` commands, we can use Bazel to run the `go` SDK.

Doing so ensures that all workstations use the same version of Go.
```bash
bazel run @rules_go//go -- mod tidy -v
bazel run @rules_go//go -- get buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go
```

Refer to [these docs](https://github.com/bazelbuild/rules_go/blob/master/docs/go/core/bzlmod.md#using-a-go-sdk) for more information.

## Managing Go modules and the listed dependencies

External Go dependencies are managed by the `go_deps` module extension which is provided by [Gazelle](./gazelle.md).

After managing the `go.mod` file with `rules_go`'s GO SDK as described above, `go_deps` can parse this `go.mod` file and
perform Minimal Version Selection on all transitive Go dependencies. Using `go.mod` allows non-Bazel projects to be able
to use our Go modules.

Example usage of `go_deps` looks like:
```starlark
# MODULE.bazel

go_deps = use_extension("@gazelle//:extensions.bzl", "go_deps")
go_deps.from_file(go_mod = "//api:go.mod")

# All *direct* Go dependencies of the module have to be listed explicitly.
use_repo(
    go_deps,
    "com_github_gogo_protobuf",
    "com_github_golang_mock",
    "com_github_golang_protobuf",
    "org_golang_x_net",
)
```

Calling
```bash
bazel mod tidy
```

automatically updates the `use_repo` declaration.

## References and links

* [rules_go and Bzlmod](https://github.com/bazelbuild/rules_go/blob/master/docs/go/core/bzlmod.md)
