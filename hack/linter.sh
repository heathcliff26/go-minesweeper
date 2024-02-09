#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath)/.."

if [ ! -d "${HOME}/.cache" ]; then
    mkdir "${HOME}/.cache"
fi

podman run -t \
    -v "${base_dir}":/app:z \
    -v "${HOME}/.cache":/root/.cache  \
    ghcr.io/heathcliff26/go-fyne-ci:latest \
    golangci-lint run -v --timeout 300s
