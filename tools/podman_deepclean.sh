#!/usr/bin/env bash

podman container kill --all
podman container rm --all
podman image rm --all
podman volume rm --all
