# Gazelle

Gazelle is a tool for generating `BUILD.bazel` files and targets in a Bazel project. Gazelle can also be used to manage
external Go module dependencies thanks to its native support for Go rule sets.

## Typical workflow for adding Go tools dependencies

Some Go modules are only relied on as tools (meaning they're not used directly in code but are required at runtime).

One example is `buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate` which is imported in generated
`.pb.go` files that import `buf/validate/validate.proto`.

For example, in `helloworld.pb.go` (generated from [helloworld.proto](../proto/helloworld/v1/helloworld.proto)), the
imports look like:

```go
import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)
```

To generate a Bazel target that can be depended on, do the following:
1. Add the dependency to [`tools.go`](../api/tools/tools.go)
2. Run `bazel mod tidy`
3. Run `bazel run @rules_go//go -- mod tidy`
4. Run `bazel run //:gazelle` to generate a `BUILD.bazel` file

See more in [these instructions](https://github.com/bazelbuild/rules_go/blob/master/docs/go/core/bzlmod.md#depending-on-tools).

## Known bugs/issues

Running Gazelle with `bazel run //:gazelle` can print debug messages like ([Issue #1409](https://github.com/bazelbuild/bazel-gazelle/issues/1409)):
```bash
gazelle: /Users/jeremymuhia/development/fjarm_go_grpc_bazel/BUILD.bazel: unknown directive: gazelle:prefix
```

This is the result of not including Gazelle supported language rules. The fix is including the target language in the
`gazelle_binary` definition:

```starlark
gazelle_binary(
    name = "gazelle-buf",
    languages = [
        "@gazelle//language/proto",  # Built-in rule from gazelle for Protos.
        "@rules_buf//gazelle/buf:buf",  # Loads the Buf extension.
        "@gazelle//language/go",  # Built-in rule from gazelle for Golang.
    ],
    visibility = ["//visibility:public"],
)
```

## References and links

* [Gazelle docs](https://github.com/bazelbuild/bazel-gazelle)