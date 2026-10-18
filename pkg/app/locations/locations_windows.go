//go:build windows

package locations

import (
	"os"
	"path/filepath"
)

func loadSaveFolderLocation() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, appName, "saves"), nil
}

func loadSettingsFileLocation() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, appName, settingsFilename), nil
}
