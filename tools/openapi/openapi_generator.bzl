"""
#CURRENT_TOOLCHAIN = "@bazel_tools//tools/jdk:current_java_runtime"
#COMPILE_TOOLCHAIN = "@bazel_tools//tools/jdk:toolchain_type"
#    print("$$$ cli_jar_file:", cli_jar_file.path)
#    print("$$$ config_file:", config_file.path)
#    args.add_joined(["-jar", cli_jar_file], join_with = " ")
#    java_compile_toolchain = ctx.toolchains["@bazel_tools//tools/jdk:toolchain_type"].java
#    java_runtime_toolchain = ctx.toolchains["@bazel_tools//tools/jdk:runtime_toolchain_type"].java_runtime
#    print("$$$ runtime_toolchain:", java_runtime_toolchain.java_executable_exec_path)
#    print("$$$ source_files:", [f.path for f in source_files])
#        print("$$$ arg, val:", arg, val)
#        print("$$$ file:", "-i", file.path)
#        executable = ctx.files._jdk,
#        executable = java_runtime_toolchain.java_executable_exec_path,
"""

RUNTIME_TOOLCHAIN = "@bazel_tools//tools/jdk:runtime_toolchain_type"

def _openapi_generate(ctx):
    cli_jar_file = ctx.file.cli_jar
    config_file = ctx.file.config

    source_files = []
    for source in ctx.files.specs:
        source_files.append(source)

    args = ctx.actions.args()
    args = args.set_param_file_format("shell")

    args.add("-jar")
    args.add(cli_jar_file)

    args.add("generate")
    args.add("-c")
    args.add(config_file)

    for arg, val in ctx.attr.arguments.items():
        args.add(arg)
        args.add(val)

    for file in source_files:
        args.add("-i")
        args.add(file)

    output_directory = ctx.actions.declare_directory("generated")

    ctx.actions.run(
        executable = ctx.attr._jdk[java_common.JavaRuntimeInfo].java_executable_exec_path,
        outputs = [output_directory],
        inputs = depset([cli_jar_file, config_file] + source_files + ctx.files._jdk),
        arguments = [args],
    )

    return [
        DefaultInfo(
            files = depset([
                cli_jar_file,
                config_file,
                output_directory,
            ] + source_files + ctx.files._jdk),
        ),
    ]

openapi_generate = rule(
    doc = "Run the openapi-generator-cli",
    implementation = _openapi_generate,
    attrs = {
        "cli_jar": attr.label(
            default = Label("//tools/openapi:openapi-generator-cli.jar"),
            allow_single_file = True,
        ),
        "arguments": attr.string_dict(
            mandatory = False,
            doc = "Extra args to be passed to openapi-generate",
        ),
        "config": attr.label(
            allow_single_file = True,
            mandatory = False,
            doc = "YAML containing the file to pass into the -c argument of openapi-generate",
        ),
        "specs": attr.label_list(
            allow_files = [".yaml"],
            mandatory = False,
            doc = "YAML files containing OpenAPI specs",
        ),
        "_jdk": attr.label(
            default = Label("@bazel_tools//tools/jdk:current_java_runtime"),
            providers = [java_common.JavaRuntimeInfo],
            cfg = "host",
        ),
    },
    toolchains = [
        RUNTIME_TOOLCHAIN,
    ],
)
