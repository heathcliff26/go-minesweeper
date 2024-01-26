package app

import (
	"os"
	"runtime"
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	a := test.NewApp()
	v := getVersion(a)

	assert := assert.New(t)

	if a.Metadata().Name != "" {
		assert.Contains(a.Metadata().Name, v.Name)
	} else {
		assert.Contains(os.Args[0], v.Name)
	}
	assert.Equal("v"+a.Metadata().Version, v.Version)
	assert.LessOrEqual(len(v.Commit), 7)
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
