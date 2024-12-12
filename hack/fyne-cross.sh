#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath)/.."

FYNE_CROSS_IMAGE="ghcr.io/heathcliff26/go-fyne-ci:latest"
APP_ID="io.github.heathcliff26.go-minesweeper"

os="${1}"
arches="${2:-$(go env GOARCH)}"
fyne_cross="$(go env GOPATH)/bin/fyne-cross"

if [ ! -f "${fyne_cross}" ]; then
    go install github.com/fyne-io/fyne-cross@latest
fi

pushd "${base_dir}" >/dev/null

${fyne_cross} "${os}" -arch="${arches}" -image="${FYNE_CROSS_IMAGE}" ./cmd/app/

IFS=',' read -ra arch_array <<<"${arches}"

for arch in "${arch_array[@]}"; do
    if [ "${os}" == "linux" ]; then
        mv "fyne-cross/bin/linux-${arch}/app" "fyne-cross/bin/linux-${arch}/go-minesweeper"
        rm -rf "fyne-cross/dist/linux-${arch}"

        echo "Building actual package fyne-cross/dist/go-minesweeper_linux-${arch}.tar.gz"
        tmp_dir="fyne-cross/tmp/linux-${arch}-packaging"
        [ -e "${tmp_dir}" ] && rm -rf "${tmp_dir}"
        mkdir "${tmp_dir}"
        cp packages/* "fyne-cross/bin/linux-${arch}/go-minesweeper" "${tmp_dir}/"
        cp img/mine.png "${tmp_dir}/${APP_ID}.png"
        tar -C "${tmp_dir}" -czf "fyne-cross/dist/go-minesweeper_linux-${arch}.tar.gz" .
        rm -rf "${tmp_dir}"
    elif [ "${os}" == "windows" ]; then
        mv "fyne-cross/dist/windows-${arch}/go-minesweeper.exe.zip" "fyne-cross/dist/go-minesweeper_windows-${arch}.zip"
        rm -rf "fyne-cross/dist/windows-${arch}"
    fi
done

popd >/dev/null
