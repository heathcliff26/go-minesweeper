#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath | xargs dirname)"
name="$(yq -r '.project_name' "${base_dir}/.goreleaser.yaml")"

export BUILDER_IMAGE="${BUILDER_IMAGE:-ghcr.io/heathcliff26/go-fyne-ci:latest}"

echo "Building releaser artifacts with goreleaser"
[ -e "$HOME/.cache/go-build" ] || mkdir -p "$HOME/.cache/go-build"
podman run --name "${name}-builder" --rm \
    -v "${base_dir}:/app:z" \
    -v "${HOME}/.cache/go-build:/root/.cache/go-build:z" \
    "${BUILDER_IMAGE}" \
    goreleaser release --skip=announce,publish,validate --clean -p 1
