#!/usr/bin/env bash

bazel coverage --combined_report=lcov //...
genhtml --branch-coverage --output genhtml "$(bazel info output_path)/_coverage/_coverage_report.dat"

if command -v xdg-open > /dev/null; then
    xdg-open genhtml/index.html
elif command -v open > /dev/null; then
    open genhtml/index.html
elif [ -n "$BROWSER" ]; then
   $BROWSER genhtml/index.html
else
   echo "Error: Could not detect a method to open the report. Please open genhtml/index.html manually."
   exit 1
fi
