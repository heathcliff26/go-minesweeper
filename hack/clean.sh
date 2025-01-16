#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath)/.."

folders=("fyne-cross" "bin" "coverprofiles" "dist" "saves")
files=("settings.yaml" "coverprofile.out")

for folder in "${folders[@]}"; do
    if ! [ -e "${base_dir}/${folder}" ]; then
        continue
    fi
    echo "Removing ${folder}"
    rm -rf "${base_dir:-.}/${folder}"
done

for file in "${files[@]}"; do
    if ! [ -e "${base_dir}/${file}" ]; then
        continue
    fi
    echo "Removing ${file}"
    rm "${base_dir:-.}/${file}"
done
