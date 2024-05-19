package app

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/heathcliff26/go-minesweeper/pkg/minesweeper"
	"github.com/stretchr/testify/assert"
)

func TestNewGrid(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)

	assert := assert.New(t)

	assert.Equal(g.Row(), len(g.Tiles))
	assert.Equal(g.Col(), len(g.Tiles[0]))
	assert.Equal(DEFAULT_DIFFICULTY, g.Difficulty)
	assert.Nil(g.Game)
	assert.NotNil(g.Timer)
	assert.NotNil(g.MineCount)
	assert.NotNil(g.ResetButton)
}

func TestTappedTile(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
	for _, row := range g.Tiles {
		for _, tile := range row {
			tile.CreateRenderer()
		}
	}

	p := minesweeper.NewPos(0, 0)

	assert := assert.New(t)

	g.TappedTile(p)
	assert.NotNil(g.Game)
	assert.True(g.Tiles[0][0].Field.Checked)

	game := g.Game.(*minesweeper.LocalGame)

	game.GameWon = true
	game.UpdateStatus()
	g.TappedTile(p)
	assert.Equal(ResetGameWonText, g.ResetButton.Label.Text)
	game.GameWon = false

	game.GameOver = true
	game.UpdateStatus()
	g.TappedTile(p)
	assert.Equal(ResetGameOverText, g.ResetButton.Label.Text)
}

func TestNewGame(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
	for _, row := range g.Tiles {
		for _, tile := range row {
			tile.CreateRenderer()
		}
	}

	g.TappedTile(minesweeper.NewPos(0, 0))
	g.NewGame()

	assert := assert.New(t)

	for x := 0; x < g.Row(); x++ {
		for y := 0; y < g.Col(); y++ {
			assert.Falsef(g.Tiles[x][y].Field.Checked, "(%d, %d) All fields should be reset", x, y)
		}
	}

	assert.Nil(g.Game)
	assert.Equal(g.Difficulty.Mines, g.MineCount.Count)
	assert.False(g.Timer.running)
	assert.Equal(ResetDefaultText, g.ResetButton.Label.Text)
}

func TestReplay(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
	for _, row := range g.Tiles {
		for _, tile := range row {
			tile.CreateRenderer()
		}
	}

	g.TappedTile(minesweeper.NewPos(0, 0))
	g.Replay()

	assert := assert.New(t)

	for x := 0; x < g.Row(); x++ {
		for y := 0; y < g.Col(); y++ {
			assert.Falsef(g.Tiles[x][y].Field.Checked, "(%d, %d) All fields should be reset", x, y)
		}
	}

	assert.NotNil(g.Game)
	assert.True(g.Game.IsReplay())
	assert.Equal(g.Difficulty.Mines, g.MineCount.Count)
	assert.False(g.Timer.running)
	assert.Equal(ResetDefaultText, g.ResetButton.Label.Text)
}

func TestUpdateFromStatus(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
	// Should not panic
	g.updateFromStatus(nil)
}

func TestAssistedMode(t *testing.T) {
	assert := assert.New(t)

	save, err := minesweeper.LoadSave("../minesweeper/testdata/assisted_mode_1.sav")
	if !assert.Nil(err, "Should load savegame") {
		t.FailNow()
	}
	g := NewMinesweeperGrid(save.Data.Difficulty, true)
	for _, row := range g.Tiles {
		for _, tile := range row {
			tile.CreateRenderer()
		}
	}
	g.Game = save.Game()

	buf, err := os.ReadFile("../minesweeper/testdata/assisted_mode_1.json")
	if !assert.Nil(err, "Should load test config") {
		t.FailNow()
	}

	var testConfig struct {
		StartPos minesweeper.Pos
		Mines    []minesweeper.Pos
		SafePos  []minesweeper.Pos
	}
	err = json.Unmarshal(buf, &testConfig)
	if !assert.Nil(err, "Should parse test config") {
		t.FailNow()
	}

	g.TappedTile(testConfig.StartPos)

	mines := 0
	safePos := 0
	for x := 0; x < g.Row(); x++ {
		for y := 0; y < g.Col(); y++ {
			p := minesweeper.NewPos(x, y)
			switch g.Tiles[p.X][p.Y].Marker {
			case HelpMarkingNone:
				assert.NotContains(testConfig.Mines, p, "Should have been marked as a mines")
				assert.NotContains(testConfig.SafePos, p, "Should have been marked as safe")
			case HelpMarkingMine:
				mines++
				assert.Contains(testConfig.Mines, p, "Should not have been marked as a mine")
			case HelpMarkingSafe:
				safePos++
				assert.Contains(testConfig.SafePos, p, "Should not have been marked as safe")
			default:
				assert.Fail("Found unknown Marker %d at %s", g.Tiles[x][y].Marker, p.String())
			}
		}
	}
	assert.Equal(len(testConfig.Mines), mines, "Should have the same amount of mines as in the config")
	assert.Equal(len(testConfig.SafePos), safePos, "Should have the same amount of safe positions as in the config")
}
