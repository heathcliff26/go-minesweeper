//go:build linux && !android

package locations

import (
	"os"
	"path/filepath"
)

func loadSaveFolderLocation() (string, error) {
	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome != "" {
		return filepath.Join(xdgDataHome, "saves"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", appName, "saves"), nil
}

func loadSettingsFileLocation() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, appName, settingsFilename), nil
}
