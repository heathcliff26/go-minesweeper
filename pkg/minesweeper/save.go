package minesweeper

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
)

const SaveFileExtension = ".sav"

type Save struct {
	// ID is the hash generated from the Data
	ID string `json:"id"`
	// Data contains everything necessary to create a game
	Data saveData `json:"data"`
}

type saveData struct {
	Mines      []Pos      `json:"mines"`
	Difficulty Difficulty `json:"difficulty"`
}

// Create a new save from the given game
func NewSave(game *LocalGame) (*Save, error) {
	data := saveData{
		Mines:      game.getMines(),
		Difficulty: game.Difficulty(),
	}

	buf, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	hash := sha512.Sum512(buf)

	return &Save{
		ID:   hex.EncodeToString(hash[:]),
		Data: data,
	}, nil
}

// Load a save file from the given path
func LoadSave(path string) (*Save, error) {
	// #nosec G304 -- Local users can decide on their file path themselves.
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var save Save
	err = json.Unmarshal(buf, &save)
	if err != nil {
		return nil, err
	}

	return &save, nil
}

// Write save file to the given path.
// Needs to have the correct extension.
func (s *Save) Save(path string) error {
	if filepath.Ext(path) != SaveFileExtension {
		path += SaveFileExtension
	}

	data, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return err
	}

	// #nosec G306 -- File permission should be determined by umask.
	return os.WriteFile(path, data, 0644)
}

// Creates a LocalGame from the save.
// The game will be considered a replay
func (s *Save) Game() *LocalGame {
	g := blankGame(s.Data.Difficulty)

	for _, p := range s.Data.Mines {
		g.Field[p.X][p.Y].Content = Mine
	}

	g.calculateFieldContent()

	g.replay = true

	return g
}
