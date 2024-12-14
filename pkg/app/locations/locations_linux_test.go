//go:build linux

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
	assert.Contains(path, filepath.Join(".config", appName, "saves"), "Variable should have home location ending")

	t.Setenv("XDG_DATA_HOME", "test")
	path, err = loadSaveFolderLocation()
	assert.NoError(err)
	assert.Equal("test/saves", path, "Should read data location from env variable")
}

func TestLoadSettingsFileLocation(t *testing.T) {
	assert := assert.New(t)

	path, err := loadSettingsFileLocation()
	assert.NoError(err)
	assert.Contains(path, filepath.Join(".config", appName, settingsFilename), "Variable should have home location ending")

	t.Setenv("XDG_CONFIG_HOME", "test")
	path, err = loadSettingsFileLocation()
	assert.NoError(err)
	assert.Equal("test/"+settingsFilename, path, "Should read settings.yaml location from env variable")
}
