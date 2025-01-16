#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath)/.."

pushd "${base_dir}" >/dev/null

OUT_DIR="coverprofiles"

if [ ! -d "${OUT_DIR}" ]; then
    mkdir "${OUT_DIR}"
fi

"${base_dir}/hack/unit-tests.sh"
go tool cover -html "coverprofile.out" -o "${OUT_DIR}/index.html"

popd >/dev/null
