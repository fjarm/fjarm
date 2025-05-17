#!/usr/bin/env bash

bazel coverage --combined_report=lcov //...
genhtml --branch-coverage --output genhtml "$(bazel info output_path)/_coverage/_coverage_report.dat"
open genhtml/index.html
