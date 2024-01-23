#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath)/.."

declare -a ALL_OS=("windows" "linux")

fyne_cross="$(go env GOPATH)/bin/fyne-cross"

if [ ! -f "${fyne_cross}" ]; then
    go install github.com/fyne-io/fyne-cross@latest
fi

for os in "${ALL_OS[@]}"; do
    "${base_dir}/hack/fyne-cross.sh" "${os}" amd64,arm64
done
