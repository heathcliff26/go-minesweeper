package app

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/heathcliff26/go-minesweeper/pkg/minesweeper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	opts := slog.HandlerOptions{
		Level: slog.LevelError,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &opts))
	slog.SetDefault(logger)
}

func TestApp(t *testing.T) {
	overrideSettingsPath = "not-a-file.yaml"
	t.Cleanup(func() {
		overrideSettingsPath = ""
	})

	newApp = test.NewApp

	a := New()

	t.Run("App", func(t *testing.T) {
		assert.NotNil(t, a.app)
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
		assert.Equal(minesweeper.Difficulties()[DEFAULT_DIFFICULTY], a.grid.Difficulty)
		assert.False(a.grid.AssistedMode)
		assert.Equal(DEFAULT_GAME_ALGORITHM, a.grid.GameAlgorithm)
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
		for _, opt := range a.gameMenu[:2] {
			t.Run(opt.Label, func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				a.grid.TappedTile(minesweeper.NewPos(0, 0))
				require.True(a.grid.Timer.Running(), "The timer should be running")

				opt.Action()
				assert.False(a.grid.Timer.Running(), "Game should not be running")

				a.NewGrid(minesweeper.Difficulties()[DEFAULT_DIFFICULTY])
				a.grid.TappedTile(minesweeper.NewPos(0, 0))
				require.True(a.grid.Timer.Running(), "The timer should be running")

				opt.Action()
				assert.False(a.grid.Timer.Running(), "Game should not be running")
			})
		}
	})
	t.Run("AssistedMode", func(t *testing.T) {
		a.NewGrid(minesweeper.Difficulties()[DEFAULT_DIFFICULTY])

		assert := assert.New(t)

		assert.False(a.assistedMode.Checked)
		assert.Equal(a.assistedMode.Checked, a.grid.AssistedMode)

		a.assistedMode.Action()
		assert.True(a.assistedMode.Checked)
		assert.Equal(a.assistedMode.Checked, a.grid.AssistedMode)

		a.assistedMode.Action()
		assert.False(a.assistedMode.Checked)
		assert.Equal(a.assistedMode.Checked, a.grid.AssistedMode)
	})
	t.Run("GameAlgorithm", func(t *testing.T) {
		for id, algorithm := range a.gameAlgorithms {
			t.Run(algorithm.Label, func(t *testing.T) {
				algorithm.Action()

				assert := assert.New(t)

				assert.Equal(id, a.grid.GameAlgorithm, "Should have set grids Algorithm")
				assert.True(algorithm.Checked, "Current Algorithm should be checked")
				for _, item := range a.gameAlgorithms {
					if item == algorithm {
						continue
					}
					assert.False(item.Checked, "No other algorithm should be checked")
				}
			})
		}
	})
}

func TestAppWithPreferences(t *testing.T) {
	overrideSettingsPath = "testdata/settings.yaml"
	t.Cleanup(func() {
		overrideSettingsPath = ""
	})

	newApp = test.NewApp

	a := New()

	p := Preferences{
		DifficultyInt: minesweeper.DifficultyClassic,
		AssistedMode:  true,
		GameAlgorithm: GameAlgorithmSafePos,
	}

	assert := assert.New(t)

	for i := range a.difficulties {
		assert.Equal(i == p.DifficultyInt, a.difficulties[i].Checked, "Only the difficulty from preferences should be selected")
	}
	for i := range a.gameAlgorithms {
		assert.Equal(i == p.GameAlgorithm, a.gameAlgorithms[i].Checked, "Only the game algorithm from preferences should be selected")
	}
	assert.Equal(p.GameAlgorithm, a.grid.GameAlgorithm, "The grid should have the correct game algorithm selected")
	assert.Equal(p.AssistedMode, a.assistedMode.Checked, "Assisted Mode should be selected")
}

func TestSaveGame(t *testing.T) {
	overrideSettingsPath = "not-a-file.yaml"
	t.Cleanup(func() {
		overrideSettingsPath = ""
	})

	newApp = test.NewApp

	a := New()
	a.grid.TappedTile(minesweeper.NewPos(0, 0))

	dir := t.TempDir()

	t.Run("NewSave", func(t *testing.T) {
		assert := assert.New(t)

		savePath := filepath.Join(dir, "new-save.sav")

		a.saveGameCallback(savePath, nil)

		assert.FileExists(savePath, "Save file should exist")
		_, err := minesweeper.LoadSave(savePath)
		assert.NoError(err, "Should be able to load the save file")
	})

	t.Run("OverwriteSave", func(t *testing.T) {
		assert := assert.New(t)

		savePath := filepath.Join(dir, "overwrite-save.sav")

		require.NoError(t, os.WriteFile(savePath, []byte("test"), 0644), "Should create an empty file")

		a.saveGameCallback(savePath, nil)

		assert.FileExists(savePath, "Save file should exist")
		_, err := minesweeper.LoadSave(savePath)
		assert.NoError(err, "Should be able to load the save file")
	})
}

func TestLoadSave(t *testing.T) {
	overrideSettingsPath = "not-a-file.yaml"
	t.Cleanup(func() {
		overrideSettingsPath = ""
	})

	newApp = test.NewApp

	a := New()

	a.loadSaveCallback("testdata/hint.sav", nil)
	assert.NotNil(t, a.grid.Game, "Game should be loaded")
}
