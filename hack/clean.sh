#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath)/.."

folders=("bin" "coverprofiles" "dist" "saves" "tmp" "x86_64" "aarch64" ".idea")
files=("settings.yaml" "coverprofile.out" "debug.keystore" "go-minesweeper.iml" "FyneApp.toml" "fyne_metadata_init.go")

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
