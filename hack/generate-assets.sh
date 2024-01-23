#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath)/.."

fyne="$(go env GOPATH)/bin/fyne-cross"

if [ ! -f "${fyne}" ]; then
    go install fyne.io/fyne/v2/cmd/fyne@latest
fi

pushd "${base_dir}" >/dev/null
go generate -v ./...
popd >/dev/null
