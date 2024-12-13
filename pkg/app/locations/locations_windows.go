//go:build windows

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
	home, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Failed to get user home folder", "err", err)
		return
	}

	saveFolder = filepath.Join(home, "AppData", "Roaming", appName, "saves")
}
