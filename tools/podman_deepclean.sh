#!/usr/bin/env bash

kind delete cluster
podman container kill --all
podman container rm --all
podman image rm --all
podman volume rm --all
