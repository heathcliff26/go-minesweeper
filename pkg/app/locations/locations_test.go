package locations

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveFolder(t *testing.T) {
	assert := assert.New(t)

	assert.Empty(saveFolder, "Should not initialize variable")

	t.Cleanup(func() {
		saveFolder = ""
	})

	path, err := SaveFolder()
	assert.NoError(err)
	assert.Equal(path, saveFolder, "Should cache path in global variable")
	assert.NotEmpty(path, "Should return a path")
	_, err = os.Stat(path)
	assert.NoError(err, "Folder should already exist or be created")

	saveFolder = "not-the-value"
	path, err = SaveFolder()
	assert.NoError(err)
	assert.Equal(path, saveFolder, "result and global variable should be equal")
	assert.Equal("not-the-value", path, "Should just return the global variable")
}

func TestSettingsFile(t *testing.T) {
	assert := assert.New(t)

	assert.Empty(settingsFile, "Should not initialize variable")

	t.Cleanup(func() {
		settingsFile = ""
	})

	path, err := SettingsFile()
	assert.NoError(err)
	assert.Equal(path, settingsFile, "Should cache path in global variable")
	assert.NotEmpty(path, "Should return a path")
	_, err = os.Stat(filepath.Dir(path))
	assert.NoError(err, "Folder should already exist or be created")

	settingsFile = "not-the-file"
	path, err = SettingsFile()
	assert.NoError(err)
	assert.Equal(path, settingsFile, "result and global variable should be equal")
	assert.Equal("not-the-file", path, "Should just return the global variable")
}
