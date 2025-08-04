package app

import (
	"os"
	"runtime"
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	oldGitCommit := gitCommit
	defer func() { gitCommit = oldGitCommit }()

	gitCommit = "1234567890abcdef"

	a := test.NewApp()
	v := getVersion(a)

	assert := assert.New(t)

	if a.Metadata().Name != "" {
		assert.Contains(a.Metadata().Name, v.Name)
	} else {
		assert.Contains(os.Args[0], v.Name)
	}
	assert.Equal("v"+a.Metadata().Version, v.Version)
	assert.Equal("1234567", v.Commit, "commit hash should be truncated")
	assert.Equal(runtime.Version(), v.Go)
}

func TestNewGridLabel(t *testing.T) {
	l := newGridLabel("test")

	assert := assert.New(t)

	assert.Equal("test", l.Text)
	assert.Equal(GridLabelColor, l.Color)
	assert.Equal(GridLabelSize, l.TextSize)
	assert.True(l.TextStyle.Bold)
}

func TestInitGitCommit(t *testing.T) {
	oldGitCommit := gitCommit
	defer func() { gitCommit = oldGitCommit }()
	assert := assert.New(t)

	gitCommit = "1234567890abcdef"
	initGitCommit()
	assert.Equal("1234567890abcdef", gitCommit, "gitCommit should not be changed")

	gitCommit = "$Format:%H$"
	initGitCommit()
	assert.NotEqual("$Format:%H$", gitCommit, "gitCommit should be changed")
}
