#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath)/.."

os="${1}"
arches="${2:-$(go env GOARCH)}"

fyne_cross="$(go env GOPATH)/bin/fyne-cross"

if [ ! -f "${fyne_cross}" ]; then
    go install github.com/fyne-io/fyne-cross@latest
fi

pushd "${base_dir}" >/dev/null

${fyne_cross} "${os}" -arch="${arches}" ./cmd/app/

IFS=',' read -ra arch_array <<<"${arches}"

for arch in "${arch_array[@]}"; do
    if [ "${os}" == "linux" ]; then
        mv "fyne-cross/bin/linux-${arch}/app" "fyne-cross/bin/linux-${arch}/go-minesweeper"
        rm -rf "fyne-cross/dist/linux-${arch}"
        tar -C "fyne-cross/bin/linux-${arch}" -czf "fyne-cross/dist/go-minesweeper_linux-${arch}.tar.gz" go-minesweeper
    elif [ "${os}" == "windows" ]; then
        mv "fyne-cross/dist/windows-${arch}/go-minesweeper.exe.zip" "fyne-cross/dist/go-minesweeper_windows-${arch}.zip"
        rm -rf "fyne-cross/dist/windows-${arch}"
    fi
done

popd >/dev/null
