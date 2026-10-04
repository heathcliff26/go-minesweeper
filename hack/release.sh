#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath | xargs dirname)"
dist_dir="${base_dir}/dist"

echo "Building releaser artifacts with goreleaser"
goreleaser release --skip=announce,publish,validate --clean -p 1

echo "Cleaning up dist directory"
rm -r "${dist_dir}/artifacts.json" "${dist_dir}/config.yaml" "${dist_dir}/metadata.json" "${dist_dir}"/*_*_checksums.txt "${dist_dir}"/linux_linux_* "${dist_dir}"/windows_windows_*
