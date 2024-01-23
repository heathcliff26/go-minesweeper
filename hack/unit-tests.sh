#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath)/.."

podman run -t -v "${base_dir}":/app:z ghcr.io/heathcliff26/go-fyne-ci:latest go test -v ./...
