#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath | xargs dirname)"

cd "${base_dir}"

sed -i 's#ID = "io.github.heathcliff26.go-minesweeper"#ID = "io.github.heathcliff26.gominesweeper"#g' FyneApp.toml

if [ -z "${KEYSTORE}" ]; then
    echo "No keystore specified, using debug keystore"
    export KEYSTORE="debug.keystore"
    export KEYSTORE_PASS="android"
    export KEYSTORE_ALIAS="androiddebugkey"
    if [ ! -e "${KEYSTORE}" ]; then
        keytool -genkeypair -v \
            -keystore "${KEYSTORE}" \
            -storepass "${KEYSTORE_PASS}" \
            -alias "${KEYSTORE_ALIAS}" \
            -keyalg RSA \
            -keysize 2048 \
            -validity 10000 \
            -dname "CN=go-minesweeper Debug,O=go-minesweeper,C=DE"
    fi
fi

fyne release --os android --app-build 1 \
    --keystore "${KEYSTORE}" \
    --keystore-pass "${KEYSTORE_PASS}" \
    --key-name "${KEYSTORE_ALIAS}"

mkdir -p dist
mv go_minesweeper.aab dist/

targets=("universal" "arm" "arm64" "amd64")

for target in "${targets[@]}"; do
    os="android"
    if [ "${target}" != "universal" ]; then
        os="android/${target}"
    fi
    echo "Building for ${target}"
    fyne package --os "${os}" --release --app-build 1
    apksigner sign \
        --ks "${KEYSTORE}" \
        --ks-pass "pass:${KEYSTORE_PASS}" \
        --ks-key-alias "${KEYSTORE_ALIAS}" \
        --v1-signing-enabled true \
        --v2-signing-enabled true \
        --v3-signing-enabled true \
        go_minesweeper.apk
    mv go_minesweeper.apk dist/go-minesweeper-"${target}".apk
done

rm go_minesweeper.apk.idsig
sed -i 's#ID = "io.github.heathcliff26.gominesweeper"#ID = "io.github.heathcliff26.go-minesweeper"#g' FyneApp.toml
sed -i 's#Build = .*#Build = 0#g' FyneApp.toml
