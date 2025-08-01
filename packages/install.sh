#!/bin/bash

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath)"

APP_ID="io.github.heathcliff26.go-minesweeper"
BINARY="go-minesweeper"

bin_dir="$HOME/.local/bin"
icon_dir="$HOME/.local/share/icons/hicolor"
if [ "$(whoami)" == "root" ]; then
    bin_dir="/usr/local/bin"
    icon_dir="/usr/share/icons/hicolor"
fi

help() {
    echo "Integrate go-minesweeper with common desktop environments."
    echo
    echo "Usage: -i | --install    -- install desktop file"
    echo "       -u | --uninstall  -- uninstall desktop file"
    echo "       -h | --help       -- show usage"
}

install() {
    echo "Installing binary to ${bin_dir}/${BINARY}"
    cp "${base_dir}/${BINARY}" "${bin_dir}/${BINARY}"

    echo "Installing desktop file"
    xdg-desktop-menu install "${base_dir}/${APP_ID}.desktop"

    echo "Installing icon"
    xdg-icon-resource install --size 512 "${base_dir}/${APP_ID}.png"
    mkdir -p "${icon_dir}/scalable/apps"
    cp "${base_dir}/${APP_ID}.svg" "${icon_dir}/scalable/apps/${APP_ID}.svg"

    xdg-desktop-menu forceupdate
    xdg-icon-resource forceupdate
}

uninstall() {
    echo "Removing binary"
    rm "${bin_dir}/${BINARY}"

    echo "Removing desktop file and icon"
    xdg-desktop-menu uninstall "${APP_ID}.desktop"
    xdg-icon-resource uninstall --size 512 "${APP_ID}.png"
    rm "${icon_dir}/scalable/apps/${APP_ID}.svg"
}

while [[ "$#" -gt 0 ]]; do
    case $1 in
    -i | --install)
        install
        exit 0
        ;;
    -u | --uninstall)
        uninstall
        exit 0
        ;;
    -h | --help)
        help
        exit 0
        ;;
    *)
        echo "Unknown argument: $1"
        help
        exit 1
        ;;
    esac
    shift
done

help
