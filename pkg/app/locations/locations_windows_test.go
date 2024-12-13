//go:build windows

package locations

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadSaveFolderLocation(t *testing.T) {
	assert.NotEmpty(t, saveFolder, "Should have initialized save folder")
	assert.True(t, strings.HasSuffix(saveFolder, "\\saves"), "Should end in folder \"saves\"")
}
