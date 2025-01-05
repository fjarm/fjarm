# Gazelle

Gazelle is a tool for generating `BUILD.bazel` files and targets in a Bazel project. Gazelle can also be used to manage
external Go module dependencies thanks to its native support for Go rule sets.

## Typical workflow for adding Go tools dependencies

Some Go modules are only relied on as tools (meaning they're not used directly in code but are required at runtime).

One example is `buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate` which is imported in generated
`.pb.go` files that import `buf/validate/validate.proto`.

In `helloworld.pb.go` (generated from [helloworld.proto](../proto/helloworld/v1/helloworld.proto)), the imports look
like:

```go
package v1

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
    * Alternatively, if the module is not just a tool and used in production, the import should still be added in a
      similar way. But, the import can be added directly in the file where it's used like in [service.go](../api/internal/helloworld/service.go):
```go
package helloworld

import (
	_ "buf.build/gen/go/fjarm/fjarm/grpc/go/helloworld/v1/helloworldv1grpc"
	// ...
)
```
2. Run `bazel mod tidy`
3. **From the `api` directory**, run `bazel run @rules_go//go -- mod tidy`
4. Run `bazel mod tidy` again
5. Run `bazel run //:gazelle` to generate a `BUILD.bazel` file

See more in [these instructions](https://github.com/bazelbuild/rules_go/blob/master/docs/go/core/bzlmod.md#depending-on-tools).

## References and links

* [Gazelle docs](https://github.com/bazelbuild/bazel-gazelle)
