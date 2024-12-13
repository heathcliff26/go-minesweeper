//go:build linux

package locations

import (
	"log/slog"
	"os"
	"path/filepath"
)

func init() {
	loadSaveFolderLocation()
}

func loadSaveFolderLocation() {
	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome != "" {
		saveFolder = filepath.Join(xdgDataHome, "saves")
		return
	}

	home, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Failed to get user home folder", "err", err)
		return
	}

	saveFolder = filepath.Join(home, ".config", appName, "saves")
}
