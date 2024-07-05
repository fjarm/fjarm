# rules_oci

This doc discusses how to use Bazel and `rules_oci` to deploy service containers.

## Known bugs

`oci.pull` and `oci_push` extensions and rules require being logged into Docker via the CLI. Be sure to run
`docker login` so these targets can run correctly.

When running containers locally on macOS, uncomment the `run --platforms=@rules_go//go/toolchain:linux_amd64` line in
[.bazelrc](../.bazelrc). This is a workaround for `exec format error`'s encountered when running `linux/amd64`
containers.

## References and links

* [oci_go_image example](https://github.com/aspect-build/bazel-examples/blob/main/oci_go_image/BUILD.bazel)
