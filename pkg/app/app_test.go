package app

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/heathcliff26/go-minesweeper/pkg/minesweeper"
	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {
	newApp = test.NewApp

	a := New()

	t.Run("App", func(t *testing.T) {
		assert.NotEmpty(t, a.app)
	})
	t.Run("Main", func(t *testing.T) {
		assert := assert.New(t)

		assert.NotEmpty(a.main)
		a.setContent()
		assert.Equal(a.Version.Name, a.main.Title())
		assert.True(a.main.FixedSize(), "Window size should be fixed")
	})
	t.Run("Version", func(t *testing.T) {
		assert := assert.New(t)

		assert.NotEmpty(a.Version)
		assert.Equal(getVersion(a.app), a.Version)
	})
	t.Run("Grid", func(t *testing.T) {
		assert := assert.New(t)

		assert.NotEmpty(a.grid)
		assert.Equal(DEFAULT_DIFFICULTY, a.grid.Difficulty)
	})
	t.Run("Difficulties", func(t *testing.T) {
		assert := assert.New(t)

		assert.NotEmpty(a.difficulties)

		diffs := minesweeper.Difficulties()
		assert.Equal(len(diffs), len(a.difficulties), "Should have a menu item for each difficulty")

		for i, d := range diffs {
			assert.Equal(d.Name, a.difficulties[i].Label, "Menu item label should match difficulty name")
		}
	})
}
