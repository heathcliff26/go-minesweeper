package locations

import (
	"log/slog"
	"os"
	"path/filepath"
)

const appName = "go-minesweeper"

const settingsFilename = "settings.yaml"

var saveFolder = ""

var settingsFile = ""

func SaveFolder() (string, error) {
	if saveFolder != "" {
		return saveFolder, nil
	}

	var err error
	saveFolder, err = loadSaveFolderLocation()
	if err == nil {
		return saveFolder, os.MkdirAll(saveFolder, 0755)
	}
	slog.Error("Failed to load save folder location, falling back to current directory", "err", err)

	return os.Getwd()
}

func SettingsFile() (string, error) {
	if settingsFile != "" {
		return settingsFile, nil
	}

	var err error
	settingsFile, err = loadSettingsFileLocation()
	if err == nil {
		return settingsFile, os.MkdirAll(filepath.Dir(settingsFile), 0755)
	}
	slog.Error("Failed to load settings.yaml location, falling back to current directory", "err", err)

	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(pwd, settingsFilename), nil
}
