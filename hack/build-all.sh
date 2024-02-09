#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath)/.."

declare -a ALL_OS=("windows" "linux")

for os in "${ALL_OS[@]}"; do
    "${base_dir}/hack/fyne-cross.sh" "${os}" amd64,arm64
done
