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
		assert.Equal(len(diffs), len(a.difficulties)-2, "Should have a menu item for each difficulty, as well as the custom option (+separator)")

		for i, d := range diffs {
			assert.Equal(d.Name, a.difficulties[i].Label, "Menu item label should match difficulty name")

			a.difficulties[i].Action()

			assert.Equal(d, a.grid.Difficulty, "Should have created grid with given difficulty")

			for n, item := range a.difficulties {
				assert.Equal((i == n), item.Checked, "Only the selected difficulty should be checked")
			}
		}
	})
	t.Run("GameMenu", func(t *testing.T) {
		for _, opt := range a.gameMenu {
			t.Run(opt.Label, func(t *testing.T) {
				assert := assert.New(t)

				a.grid.TappedTile(minesweeper.NewPos(0, 0))
				if !assert.True(a.grid.Timer.Running(), "Assert that a game is running") {
					t.FailNow()
				}

				opt.Action()
				assert.False(a.grid.Timer.Running(), "Game should not be running")

				a.grid = NewMinesweeperGrid(DEFAULT_DIFFICULTY)
				a.setContent()
				a.grid.TappedTile(minesweeper.NewPos(0, 0))
				if !assert.True(a.grid.Timer.Running(), "Assert that a game is running") {
					t.FailNow()
				}

				opt.Action()
				assert.False(a.grid.Timer.Running(), "Game should not be running")
			})
		}
	})
}
