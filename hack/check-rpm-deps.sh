#!/bin/bash

set -euo pipefail

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath | xargs dirname)"

binary_path="${base_dir}/bin/go-minesweeper"

if [ ! -f "${binary_path}" ]; then
    echo "Error: Binary not found at ${binary_path}. Build it first."
    exit 1
fi

if ! command -v yq &> /dev/null; then
    echo "Error: yq is not installed. Please install yq first."
    exit 1
fi

echo "Scanning dependencies for ${binary_path}..."

# 1. Extract .so files using ldd, filter paths, and find owning RPM packages
SO_FILES=$(ldd "${binary_path}" | grep -oE '/[^ ]+' || true)

RPM_DEPS=()
for so in ${SO_FILES}; do
    if [ -f "${so}" ]; then
        pkg=$(rpm -qf "${so}" --qf "%{NAME}\n" 2>/dev/null || true)
        if [ -n "${pkg}" ] && [ "${pkg}" != "not owned" ]; then
            RPM_DEPS+=("${pkg}")
        fi
    fi
done

# Deduplicate and sort packages into a JSON array format for yq
UNIQUE_DEPS=$(printf '%s\n' "${RPM_DEPS[@]}" | sort -u | jq -R . | jq -s .)

echo "Found RPM dependencies: ${UNIQUE_DEPS}"
