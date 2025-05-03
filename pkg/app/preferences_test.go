package app

import (
	"os"
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/heathcliff26/go-minesweeper/pkg/minesweeper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadPreferences(t *testing.T) {
	assert := assert.New(t)

	overrideSettingsPath = "testdata/not-a-file.yaml"
	t.Cleanup(func() {
		overrideSettingsPath = ""
	})

	p, err := LoadPreferences()
	assert.Error(err)
	assert.Equal(defaultPreferences(), p, "Should return the default settings when none found")

	overrideSettingsPath = "testdata/settings.yaml"
	p, err = LoadPreferences()
	assert.NoError(err)
	res := Preferences{
		DifficultyInt: minesweeper.DifficultyClassic,
		AssistedMode:  true,
		GameAlgorithm: GameAlgorithmSafePos,
	}
	assert.Equal(res, p, "Should load the settings")

	overrideSettingsPath = "testdata/not-a-yaml-file.txt"
	p, err = LoadPreferences()
	assert.Error(err)
	assert.Equal(defaultPreferences(), p, "Should return the default settings on error")
}

func TestCreatePreferencesFromApp(t *testing.T) {
	newApp = test.NewApp

	a := New()

	for i := range a.difficulties {
		a.difficulties[i].Checked = i == 0
	}
	for i := range a.gameAlgorithms {
		a.gameAlgorithms[i].Checked = i == 0
	}
	a.assistedMode.Checked = true

	p := CreatePreferencesFromApp(a)

	assert := assert.New(t)

	res := Preferences{
		DifficultyInt: minesweeper.DifficultyClassic,
		AssistedMode:  true,
		GameAlgorithm: GameAlgorithmSafePos,
	}
	assert.Equal(res, p, "Should create the Preferences correctly")

	for i := range a.difficulties {
		a.difficulties[i].Checked = false
	}
	for i := range a.gameAlgorithms {
		a.gameAlgorithms[i].Checked = false
	}
	a.assistedMode.Checked = false

	p = CreatePreferencesFromApp(a)
	assert.Equal(defaultPreferences(), p, "Should use defaults if not possible")
}

func TestPreferencesDifficuly(t *testing.T) {
	assert := assert.New(t)

	p := Preferences{
		DifficultyInt: minesweeper.DifficultyClassic,
	}
	assert.Equal(minesweeper.Difficulties()[minesweeper.DifficultyClassic], p.Difficulty(), "Should return the correct difficulty")

	p.DifficultyInt = minesweeper.DifficultyExpert
	assert.Equal(minesweeper.Difficulties()[minesweeper.DifficultyExpert], p.Difficulty(), "Should return the correct difficulty")

	p.DifficultyInt = -1
	assert.Equal(minesweeper.Difficulties()[DEFAULT_DIFFICULTY], p.Difficulty(), "Should return the default difficulty if out of bounds")
	p.DifficultyInt = minesweeper.DifficultyExpert + 1
	assert.Equal(minesweeper.Difficulties()[DEFAULT_DIFFICULTY], p.Difficulty(), "Should return the default difficulty if out of bounds")
}

func TestPreferencesSave(t *testing.T) {
	p := Preferences{
		DifficultyInt: minesweeper.DifficultyClassic,
		AssistedMode:  true,
		GameAlgorithm: GameAlgorithmSafePos,
	}

	testFilePath := "test-settings.yaml"

	overrideSettingsPath = testFilePath
	t.Cleanup(func() {
		overrideSettingsPath = ""
	})

	assert := assert.New(t)

	err := p.Save()
	require.NoError(t, err, "Should save the preferences")
	t.Cleanup(func() {
		_ = os.Remove(testFilePath)
	})

	res, err := LoadPreferences()
	assert.NoError(err, "Should load the preferences without error")
	assert.Equal(p, res, "Should load the same preferences it just saved")
}
