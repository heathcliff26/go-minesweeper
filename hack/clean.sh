#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath)/.."

folders=("fyne-cross" "bin" "coverprofiles" "dist")

for folder in "${folders[@]}"; do
    if ! [ -e "${base_dir}/${folder}" ]; then
        continue
    fi
    echo "Removing ${folder}"
    rm -rf "${base_dir:-.}/${folder}"
done
