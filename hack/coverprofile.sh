#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath)/.."

pushd "${base_dir}" >/dev/null

OUT_DIR="coverprofiles"

if [ ! -d "${OUT_DIR}" ]; then
    mkdir "${OUT_DIR}"
fi

podman run -t -v "${base_dir}":/app:z ghcr.io/heathcliff26/go-fyne-ci:latest /bin/bash -c "go test -coverprofile=\"${OUT_DIR}/cover.out\" ./... && go tool cover -html \"${OUT_DIR}/cover.out\" -o \"${OUT_DIR}/index.html\""

popd >/dev/null
