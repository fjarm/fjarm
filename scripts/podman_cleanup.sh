#!/usr/bin/env sh

set -e

if ! command -v podman >/dev/null 2>&1; then
    echo "Error: podman command not found." >&2
    exit 1
fi

MACHINE_NAME="podman-machine-default"
MACHINE_STATUS="Currently running"

# Verify there's a Podman machine running
if ! podman machine list 2>/dev/null | grep -q "${MACHINE_NAME}.*${MACHINE_STATUS}"; then
    echo "Error: Podman machine '${MACHINE_NAME}' is not running."
    echo "Please start the Podman machine with: podman machine start ${MACHINE_NAME}"
    exit 1
fi

podman container kill --all
podman container rm --all
podman image rm --all
podman volume rm --all
