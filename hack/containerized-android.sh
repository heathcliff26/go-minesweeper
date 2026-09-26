#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath | xargs dirname)"

# shellcheck source=common.sh
source "${base_dir}/hack/common.sh"

export BUILDER_IMAGE_ANDROID="${BUILDER_IMAGE_ANDROID:-${BUILDER_IMAGE}-android}"
export BUILDER_IMAGE="${BUILDER_IMAGE_ANDROID}"

if [ ! -d "${HOME}/.cache" ]; then
    mkdir "${HOME}/.cache"
fi

podman run -t --rm \
    --name netrouse-builder-android \
    --env KEYSTORE="${KEYSTORE}" \
    --env KEYSTORE_PASS="${KEYSTORE_PASS}" \
    --env KEYSTORE_ALIAS="${KEYSTORE_ALIAS}" \
    -v "${base_dir}:/app:z" \
    -v "${HOME}/.cache":/root/.cache:z \
    "${BUILDER_IMAGE}" \
    hack/build-android.sh
