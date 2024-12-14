//go:build windows

package locations

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadSaveFolderLocation(t *testing.T) {
	assert := assert.New(t)
	path, err := loadSaveFolderLocation()
	assert.NoError(err)
	assert.Contains(path, filepath.Join("AppData", "Roaming", appName, "saves"))
}

func TestLoadSettingsFileLocation(t *testing.T) {
	assert := assert.New(t)
	path, err := loadSettingsFileLocation()
	assert.NoError(err)
	assert.Contains(path, filepath.Join("AppData", "Roaming", appName, settingsFilename))
}
