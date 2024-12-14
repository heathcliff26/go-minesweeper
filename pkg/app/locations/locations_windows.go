//go:build windows

package locations

import (
	"os"
	"path/filepath"
)

func loadSaveFolderLocation() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "AppData", "Roaming", appName, "saves"), nil
}

func loadSettingsFileLocation() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "AppData", "Roaming", appName, settingsFile), nil
}
