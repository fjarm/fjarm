# rules_oci

This doc discusses how to use Bazel and `rules_oci` to deploy service containers.

## Known bugs

### Registry login

`oci.pull` and `oci_push` extensions and rules require being logged into Docker via the CLI. Be sure to run
`docker login` so these targets can run correctly.

#### Fix

Logging in to a registry through the CLI is no longer needed. Publishing to the CLI is now done through [GitHub Actions](../.github/workflows/ghcr-push.yaml).

### Locally running containers

When running containers locally on macOS, uncomment the `run --platforms=@rules_go//go/toolchain:linux_amd64` line in
[.bazelrc](../.bazelrc). This is a workaround for `exec format error`'s encountered when running `linux/amd64`
containers.

#### Fix

The override is no longer needed. Instead, a `go_cross_binary` with `platform = "@rules_go//go/toolchain:linux_arm64",`
ensures that the `go_binary` target loaded into the container is always compiled for `linux/arm64`.

Find an example in the [helloworld/cmd/BUILD.bazel](../api/internal/helloworld/v1/cmd/helloworld/BUILD.bazel) file.

## References and links

* [oci_go_image example](https://github.com/aspect-build/bazel-examples/blob/main/oci_go_image/BUILD.bazel)
