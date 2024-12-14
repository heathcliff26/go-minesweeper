#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath)/.."

folders=("fyne-cross" "bin" "coverprofiles" "dist" "saves")

for folder in "${folders[@]}"; do
    if ! [ -e "${base_dir}/${folder}" ]; then
        continue
    fi
    echo "Removing ${folder}"
    rm -rf "${base_dir:-.}/${folder}"
done

if [ -e "${base_dir}/settings.yaml" ]; then
    echo "Removing settings.yaml"
    rm "${base_dir}/settings.yaml"
fi
