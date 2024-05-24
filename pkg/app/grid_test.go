package app

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/heathcliff26/go-minesweeper/pkg/minesweeper"
	"github.com/stretchr/testify/assert"
)

type assistedModeTestConfig struct {
	CheckPos minesweeper.Pos
	Mines    []minesweeper.Pos
	SafePos  []minesweeper.Pos
}

func TestNewGrid(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, true)

	assert := assert.New(t)

	assert.Equal(g.Row(), len(g.Tiles))
	assert.Equal(g.Col(), len(g.Tiles[0]))
	assert.Equal(DEFAULT_DIFFICULTY, g.Difficulty)
	assert.Nil(g.Game)
	assert.True(g.AssistedMode)
	assert.Equal(DEFAULT_GAME_ALGORITHM, g.GameAlgorithm)
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

func TestGameAlgorithm(t *testing.T) {
	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
	for _, row := range g.Tiles {
		for _, tile := range row {
			tile.CreateRenderer()
		}
	}

	tMatrix := []string{"SafePos", "safeArea"}
	for i, tCase := range tMatrix {
		t.Run(tCase, func(t *testing.T) {
			t.Cleanup(g.NewGame)
			g.GameAlgorithm = i
			g.TappedTile(minesweeper.NewPos(0, 0))
			assert.NotNil(t, g.Game)
		})
	}
	t.Run("Unknown", func(t *testing.T) {
		t.Cleanup(g.NewGame)
		g.GameAlgorithm = -1
		g.TappedTile(minesweeper.NewPos(0, 0))
		assert.Nil(t, g.Game)
	})
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

	g, testConfig, err := loadAssistedModeTest("../minesweeper/testdata/assisted_mode_1", true)
	if !assert.Nil(err, "Should have loaded the test") {
		t.Fatal(err)
	}

	for _, step := range testConfig {
		g.TappedTile(step.CheckPos)

		mines := 0
		safePos := 0
		for x := 0; x < g.Row(); x++ {
			for y := 0; y < g.Col(); y++ {
				p := minesweeper.NewPos(x, y)
				if g.Tiles[p.X][p.Y].Field.Checked {
					continue
				}
				switch g.Tiles[p.X][p.Y].Marker {
				case HelpMarkingNone:
					assert.NotContains(step.Mines, p, "Should have been marked as a mines")
					assert.NotContains(step.SafePos, p, "Should have been marked as safe")
				case HelpMarkingMine:
					mines++
					assert.Contains(step.Mines, p, "Should not have been marked as a mine")
				case HelpMarkingSafe:
					safePos++
					assert.Contains(step.SafePos, p, "Should not have been marked as safe")
				default:
					assert.Fail("Found unknown Marker %d at %s", g.Tiles[x][y].Marker, p.String())
				}
			}
		}
		assert.Equal(len(step.Mines), mines, "Should have the same amount of mines as in the config")
		assert.Equal(len(step.SafePos), safePos, "Should have the same amount of safe positions as in the config")
	}
}

func TestHint(t *testing.T) {
	assert := assert.New(t)

	g := NewMinesweeperGrid(DEFAULT_DIFFICULTY, false)
	assert.False(g.Hint(), "Should not give hint when game is nil")

	g, testConfig, err := loadAssistedModeTest("testdata/hint", false)
	if !assert.Nil(err, "Should have loaded the test") {
		t.Fatal(err)
	}
	step := testConfig[0]

	g.TappedTile(step.CheckPos)

	for _, mine := range step.Mines {
		assert.True(g.Hint(), "Should be able to display hints")
		tile := g.Tiles[mine.X][mine.Y]
		assert.Equalf(HelpMarkingMine, tile.Marker, "Tile should be marked as mine, tile=%s", tile.Pos.String())
		tile.Flagged = true
		for x := 0; x < g.Difficulty.Row; x++ {
			for y := 0; y < g.Difficulty.Col; y++ {
				tile := g.Tiles[x][y]
				if tile.Flagged || tile.Field.Checked {
					continue
				}
				if !assert.Equalf(HelpMarkingNone, tile.Marker, "No other tiles should be marked, tile=%s, mine=%s", tile.Pos.String(), mine.String()) {
					t.FailNow()
				}
			}
		}
	}
	assert.True(g.Hint(), "Should be able to display hints")
	assert.Equal(HelpMarkingSafe, g.Tiles[step.SafePos[0].X][step.SafePos[0].Y].Marker, "Tile should be marked as safe")
	for x := 0; x < g.Difficulty.Row; x++ {
		for y := 0; y < g.Difficulty.Col; y++ {
			tile := g.Tiles[x][y]
			if tile.Flagged || tile.Field.Checked || tile.Pos == step.SafePos[0] {
				continue
			}
			if !assert.Equalf(HelpMarkingNone, tile.Marker, "No other tiles should be marked, tile=%s", tile.Pos.String()) {
				t.FailNow()
			}
		}
	}

	g.TappedTile(minesweeper.NewPos(15, 6))
	assert.False(g.Hint(), "Should not be able to display hint on failed game")
}

func loadAssistedModeTest(path string, assistedMode bool) (*MinesweeperGrid, []assistedModeTestConfig, error) {
	save, err := minesweeper.LoadSave(path + ".sav")
	if err != nil {
		return nil, nil, err
	}
	g := NewMinesweeperGrid(save.Data.Difficulty, assistedMode)
	for _, row := range g.Tiles {
		for _, tile := range row {
			tile.CreateRenderer()
		}
	}
	g.Game = save.Game()

	buf, err := os.ReadFile(path + ".json")
	if err != nil {
		return nil, nil, err
	}

	var testConfig []assistedModeTestConfig
	err = json.Unmarshal(buf, &testConfig)
	if err != nil {
		return nil, nil, err
	}
	return g, testConfig, nil
}
