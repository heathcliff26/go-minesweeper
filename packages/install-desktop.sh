#!/bin/bash

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath)"

APP_ID="io.github.heathcliff26.go-minesweeper"
DESKTOP_FILE_TARGET="${HOME}/.local/share/applications/${APP_ID}.desktop"

help() {
    echo "Integrate go-minesweeper with common desktop environments."
    echo
    echo "Usage: -i | --install    -- install desktop file"
    echo "       -u | --uninstall  -- uninstall desktop file"
    echo "       -h | --help       -- show usage"
}

install() {
    sed -e "s|@BASE_DIR|${base_dir}|g" "${base_dir}/${APP_ID}.desktop" >"${DESKTOP_FILE_TARGET}"
    echo "Created file: ${DESKTOP_FILE_TARGET}"

    if command -v xdg-desktop-menu >/dev/null 2>&1; then
        echo "Forcing desktop menu update"
        xdg-desktop-menu forceupdate
        xdg-icon-resource forceupdate
    else
        echo "The app will not show up in the menu until the session is restarted"
    fi
}

uninstall() {
    rm "${DESKTOP_FILE_TARGET}"
    echo "Removed: ${DESKTOP_FILE_TARGET}"
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
