//go:build android

package locations

import (
	"path/filepath"

	"fyne.io/fyne/v2"
)

func loadSaveFolderLocation() (string, error) {
	app := fyne.CurrentApp()

	dir := app.Storage().RootURI().Path()

	return filepath.Join(dir, "saves"), nil
}

func loadSettingsFileLocation() (string, error) {
	app := fyne.CurrentApp()

	dir := app.Storage().RootURI().Path()

	return filepath.Join(dir, settingsFilename), nil
}
