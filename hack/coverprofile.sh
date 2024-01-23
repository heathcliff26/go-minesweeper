#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath)/.."

pushd "${base_dir}" >/dev/null

OUT_DIR="${base_dir}/coverprofiles"

if [ ! -d "${OUT_DIR}" ]; then
    mkdir "${OUT_DIR}"
fi

go test -coverprofile="${OUT_DIR}/cover.out" "./..."
go tool cover -html "${OUT_DIR}/cover.out" -o "${OUT_DIR}/index.html"

popd >/dev/null
