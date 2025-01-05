---

# Graveyard

## Known bugs/issues

Running Gazelle with `bazel run //api:gazelle` can print debug messages like ([Issue #1409](https://github.com/bazelbuild/bazel-gazelle/issues/1409)):
```bash
gazelle: /Users/jeremymuhia/development/fjarm/api/BUILD.bazel: unknown directive: gazelle:prefix
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

Sometimes Gazelle will write Go targets in the [proto](../proto) modules after running `bazel run //:gazelle`. This can
lead to conflicts with manually defined targets in the [api](../api) modules:

```bash
gazelle: rule //api/internal/helloworld imports "github.com/fjarm/fjarm/api/pkg/helloworld/v1" which matches multiple rules: //api/pkg/helloworld/v1:helloworld_service_proto and //api/pkg/helloworld/v1:helloworld_library. # gazelle:resolve may be used to disambiguate
```

Adding a `# gazelle:resolve` tag to the top level [BUILD.bazel](../BUILD.bazel) file should fix this.