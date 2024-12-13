//go:build linux

package locations

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadSaveFolderLocation(t *testing.T) {
	assert := assert.New(t)

	home, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Failed to get user home folder", "err", err)
		return
	}

	expectedFolder := filepath.Join(home, ".config", appName, "saves")

	assert.Equal(expectedFolder, saveFolder, "Variable should have been initialized")

	t.Setenv("XDG_DATA_HOME", "test")
	loadSaveFolderLocation()
	assert.Equal("test/saves", saveFolder, "Should read data location from env variable")
}
