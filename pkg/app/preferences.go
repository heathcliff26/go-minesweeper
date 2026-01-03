package app

import (
	"os"

	"github.com/heathcliff26/go-minesweeper/pkg/app/locations"
	"github.com/heathcliff26/go-minesweeper/pkg/minesweeper"
	"go.yaml.in/yaml/v3"
)

type Preferences struct {
	DifficultyInt int  `yaml:"difficulty"`
	AssistedMode  bool `yaml:"assistedMode"`
	GameAlgorithm int  `yaml:"gameAlgorithm"`
}

// Only used for testing, as we don't want to touch the actual file that might be on the system.
var overrideSettingsPath = ""

func defaultPreferences() Preferences {
	return Preferences{
		DifficultyInt: DEFAULT_DIFFICULTY,
		AssistedMode:  false,
		GameAlgorithm: DEFAULT_GAME_ALGORITHM,
	}
}

// Load the preferences from the settings file.
// Return default values when an error occurs.
func LoadPreferences() (Preferences, error) {
	path, err := locations.SettingsFile()
	if err != nil {
		return defaultPreferences(), err
	}
	if overrideSettingsPath != "" {
		path = overrideSettingsPath
	}

	// #nosec G304 -- File path is hardcoded relative to user directory.
	f, err := os.ReadFile(path)
	if err != nil {
		return defaultPreferences(), err
	}

	p := defaultPreferences()
	err = yaml.Unmarshal(f, &p)
	if err != nil {
		return defaultPreferences(), err
	}

	return p, nil
}

// Convert the current app state to preferences that can be saved.
func CreatePreferencesFromApp(app *App) Preferences {
	difficulty := DEFAULT_DIFFICULTY
	for i, d := range app.difficulties {
		if d.Checked {
			difficulty = i
			break
		}
	}

	gameAlgorithm := DEFAULT_GAME_ALGORITHM
	for i, algorithm := range app.gameAlgorithms {
		if algorithm.Checked {
			gameAlgorithm = i
			break
		}
	}

	return Preferences{
		DifficultyInt: difficulty,
		AssistedMode:  app.assistedMode.Checked,
		GameAlgorithm: gameAlgorithm,
	}
}

// Convert the saved difficulty number to an actual difficulty
func (p Preferences) Difficulty() minesweeper.Difficulty {
	difficulties := minesweeper.Difficulties()

	if p.DifficultyInt < 0 || p.DifficultyInt >= len(difficulties) {
		return difficulties[DEFAULT_DIFFICULTY]
	}
	return difficulties[p.DifficultyInt]
}

// Save the preferences to the settings file
func (p Preferences) Save() error {
	path, err := locations.SettingsFile()
	if err != nil {
		return err
	}
	if overrideSettingsPath != "" {
		path = overrideSettingsPath
	}

	data, err := yaml.Marshal(&p)
	if err != nil {
		return err
	}

	// #nosec G306 -- File permission should be determined by umask.
	return os.WriteFile(path, data, 0644)
}
